package audit

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// WatchAndInjectHooks finds structs with hookie.Inject and calls their hooks
func WatchAndInjectHooks(ctx context.Context) error {
	rootDir, err := findProjectRoot()
	if err != nil {
		return err
	}
	goDirs, err := collectGoDirs(rootDir)
	if err != nil {
		log.Fatalf("Error collecting Go directories: %v", err)
	}
	for _, dir := range goDirs {
		err = WatchAndRegister(ctx, dir)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// WatchAndRegister scans the specified directory for structs with in.Inject and registers their hooks.
// It returns an error if any issues occur during processing.
func WatchAndRegister(ctx context.Context, dir string) error {
	logger := slog.Default()
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedImports,
	}

	// Handle any potential panics gracefully
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic", "error", r)
		}
	}()

	// Load the packages from the specified directory
	_packages, err := packages.Load(cfg, dir)
	if err != nil {
		logger.Error("Failed to load packages from directory", "directory", dir, "error", err)
		return fmt.Errorf("unable to load packages from directory '%s': %w", dir, err) // Provide context in the returned error
	}

	// Iterate over loaded packages
	for _, pkg := range _packages {
		// Iterate over syntax trees of the package
		for _, file := range pkg.Syntax {
			// Inspect the AST for declarations
			for _, declaration := range file.Decls {
				// Check for type declarations
				if genDecl, ok := declaration.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						typeSpec := spec.(*ast.TypeSpec)
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							if isInjectable(structType) {
								structName := typeSpec.Name.Name
								logger.Info("Found injectable struct", "structName: ", structName, "packagePath", pkg.ID)

								// Lookup the full type from the package scope
								typeObject := pkg.Types.Scope().Lookup(structName)
								if typeObject == nil {
									return fmt.Errorf("struct '%s' not found in package '%s': unable to retrieve its type information", structName, pkg.PkgPath)
								}

								// Ensure the object is of type *types.Named
								if namedType, ok := typeObject.Type().(*types.Named); ok && namedType != nil {
									RegisterModel(pkg.ID + structName)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil // Return nil if processing completes without error
}

// isInjectable checks if a given struct type contains a field of type "Inject".
// It returns true if such a field is found, otherwise returns false.
func isInjectable(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		// Check if the field is a SelectorExpr (e.g., in.Inject)
		if selExpr, ok := field.Type.(*ast.SelectorExpr); ok {
			pkgIdent, ok := selExpr.X.(*ast.Ident)
			if ok && pkgIdent.Name == "in" && selExpr.Sel.Name == "Inject" {
				// Found a field of type instance.Inject
				return true
			}
		}

		// Handle qualified names (e.g., just Inject)
		if typeIdent, ok := field.Type.(*ast.Ident); ok {
			if typeIdent.Name == "Inject" {
				return true // Found a field of type Inject
			}
		}

		// Handle pointer to a type (e.g., *Inject or *instance.Inject)
		if typeSpec, ok := field.Type.(*ast.StarExpr); ok {
			// If it's a pointer to a simple type (e.g., *Inject)
			if ident, ok := typeSpec.X.(*ast.Ident); ok && ident.Name == "Inject" {
				return true // Found a pointer field of type Inject
			}
			// If it's a pointer to a qualified type (e.g., *instance.Inject)
			if selExpr, ok := typeSpec.X.(*ast.SelectorExpr); ok {
				pkgIdent, ok := selExpr.X.(*ast.Ident)
				if ok && pkgIdent.Name == "in" && selExpr.Sel.Name == "Inject" {
					return true // Found a pointer field of type instance.Inject
				}
			}
		}
	}
	return false // No field of type Inject found
}

// collectGoDirs walks through the baseDir recursively and collects directories
// that contain Go files, excluding those filtered by isExclude.
// It returns a slice of Go directories and any encountered error.
func collectGoDirs(baseDir string) ([]string, error) {
	var goDirs []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if hasGoFiles(path) && !isExclude(path) {
				goDirs = append(goDirs, path)
			}
		}
		return nil
	})
	return goDirs, err
}

// hasGoFiles checks if the specified directory contains any Go files.
func hasGoFiles(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", dir, err)
		return false
	}

	// Check each file in the directory
	for _, file := range files {
		// If a file is not a directory and has a ".go" extension, it is a Go file
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			return true
		}
	}
	return false
}

// isExclude checks if a directory should be excluded based on its prefix.
// Directories starting with "cmd" or "golang.org" are excluded.
func isExclude(dir string) bool {
	if strings.HasPrefix(dir, "cmd") ||
		strings.HasPrefix(dir, "golang.org") {
		return true
	}
	return false
}

// getGoModuleName reads and parses the "go.mod" file in the specified directory
// to retrieve the Go module name. It returns the module name or an error.
func getGoModuleName(dir string) (string, error) {
	// Construct the path to the "go.mod" file.
	goModPath := filepath.Join(dir, "go.mod")

	// Read the content of the "go.mod" file.
	data, err := os.ReadFile(goModPath)
	if err != nil {
		// Return an error if the file can't be read.
		return "", fmt.Errorf("could not read go.mod file: %v", err)
	}

	// Parse the "go.mod" file to extract the module information.
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		// Return an error if parsing the file fails.
		return "", fmt.Errorf("could not parse go.mod file: %v", err)
	}

	// Return the module name from the parsed "go.mod" file.
	return modFile.Module.Mod.Path, nil
}

// findProjectRoot searches for the nearest "go.mod" file starting from the current working directory.
// It returns the directory containing "go.mod" or an error if not found.
func findProjectRoot() (string, error) {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err // Return error if unable to get the current directory
	}

	for {
		// Construct the path to the "go.mod" file
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil // Found go.mod, return the directory
		}

		// Move up one directory level
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached the root of the filesystem
		}
		dir = parent // Update dir to the parent directory
	}

	// Return an error if no go.mod file is found in any parent directory
	return "", fmt.Errorf("go.mod file not found")
}
