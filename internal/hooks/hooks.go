package hooks

import (
	"context"
	"fmt"
	"github.com/its-own/gaudit/db/mongo"
	"github.com/its-own/gaudit/internal/audit_log"
	"github.com/its-own/gaudit/internal/entities"
	in "github.com/its-own/gaudit/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"reflect"
	"strings"
	"time"
	"unicode"
)

type DefaultHooks struct {
	l *slog.Logger
}

func NewDefaultHook() *DefaultHooks {
	return &DefaultHooks{l: slog.Default()}
}

func (h *DefaultHooks) PreSave(ctx context.Context, model interface{}, filter interface{}, col, ops, docId string) {
	// Trigger PreSave hook if defined by user, else run default
	if hasPreSaveHook(model) {
		model.(in.Hook).PreSave(ctx, model, filter, col, ops, docId)
		return
	}
	h.l.Info("default PreSave hook triggered")
}

func (h *DefaultHooks) PostSave(ctx context.Context, model interface{}, filter interface{}, col, ops, docId string) {
	if hasPostSaveHook(model) {
		model.(in.Hook).PostSave(ctx, model, filter, col, ops, docId)
		return
	}

	if isAuditLogEnabled(model) {
		switch ops {
		case "insert":
			h.handleInsertOperation(ctx, model)
		case "update":
			h.handleUpdateOperation(ctx, model, docId)
		}
	}
	h.l.Info("default PostSave hook triggered")
}

// handleInsertOperation manages audit logging during insert operations.
func (h *DefaultHooks) handleInsertOperation(ctx context.Context, model interface{}) {
	db := mongo.GetDbConnection()
	state, err := structToMap(model)
	if err != nil {
		h.l.Error(err.Error())
		return
	}

	auditLogMeta := entities.AuditLogMeta{
		Id:                   primitive.NewObjectID(),
		DocumentCurrentState: state,
	}

	_, err = db.Database.Collection("audit_logs_meta").InsertOne(ctx, auditLogMeta)
	if err != nil {
		db.Logger.Error(err.Error())
	}
}

// handleUpdateOperation manages audit logging during update operations.
func (h *DefaultHooks) handleUpdateOperation(ctx context.Context, model interface{}, docId string) {
	db := mongo.GetDbConnection()

	// Retrieve the existing audit log meta by document ID
	auditLogMeta, err := h.findAuditLogMeta(ctx, db, docId)
	if err != nil {
		h.l.Error(fmt.Sprintf("Failed to find audit log meta: %v", err))
		return
	}

	// Convert the new document state to a map
	newDoc, err := structToMap(model)
	if err != nil {
		h.l.Error(fmt.Sprintf("Failed to convert model to map: %v", err))
		return
	}

	// Compare document states and log changes
	changeLog := compareDocumentStates(auditLogMeta.DocumentCurrentState, newDoc)
	if err := h.logAuditChanges(ctx, db, changeLog); err != nil {
		h.l.Error(fmt.Sprintf("Failed to log audit changes: %v", err))
		return
	}

	// Update the audit meta with the new document state
	if err := h.updateAuditLogMeta(ctx, db, auditLogMeta.Id, newDoc); err != nil {
		h.l.Error(fmt.Sprintf("Failed to update audit log meta: %v", err))
	}
}

// Helper function to retrieve a value from the context and handle missing data.
func getContextValue(ctx context.Context, key string) string {
	value, ok := ctx.Value(key).(string)
	if !ok || value == "" {
		return "default"
	}
	return value
}

// findAuditLogMeta retrieves the existing audit log meta by document ID.
func (h *DefaultHooks) findAuditLogMeta(ctx context.Context, db *mongo.Mongo, docId string) (entities.AuditLogMeta, error) {
	var auditLogMeta entities.AuditLogMeta
	auditFilter := bson.D{{"document_current_state._id", docId}}

	err := db.Database.Collection("audit_logs_meta").FindOne(ctx, auditFilter).Decode(&auditLogMeta)
	if err != nil {
		return auditLogMeta, fmt.Errorf("error finding audit log meta: %w", err)
	}
	return auditLogMeta, nil
}

// logAuditChanges inserts a new audit log entry for document changes.
func (h *DefaultHooks) logAuditChanges(ctx context.Context, db *mongo.Mongo, changeLog map[string]entities.AuditChange) error {
	currentTime := time.Now()

	// Create the audit log entry
	auditLog := entities.AuditLog{
		Id:             primitive.NewObjectID(),
		AuditURL:       "example.com", // Example, you can replace with real URL
		AuditIPAddress: getContextValue(ctx, "ip_addr"),
		AuditUserAgent: getContextValue(ctx, "user_agent"),
		AuditTags:      []string{"audit", "log"},
		AuditCreatedAt: &currentTime,
		UserID:         getContextValue(ctx, "user_id"),
		UserType:       getContextValue(ctx, "role"),
		Change:         changeLog,
	}

	// Insert the audit log into the collection
	_, err := db.Database.Collection("audit_logs").InsertOne(context.Background(), auditLog)
	if err != nil {
		return fmt.Errorf("error inserting audit log: %w", err)
	}
	return nil
}

// updateAuditLogMeta updates the audit log meta with the new document state.
func (h *DefaultHooks) updateAuditLogMeta(ctx context.Context, db *mongo.Mongo, auditLogMetaId primitive.ObjectID, newDoc map[string]interface{}) error {
	update := bson.M{
		"$set": bson.M{"document_current_state": newDoc},
	}

	_, err := db.Database.Collection("audit_logs_meta").UpdateByID(ctx, auditLogMetaId, update)
	if err != nil {
		return fmt.Errorf("error updating audit log meta: %w", err)
	}
	return nil
}

// isAuditLogEnabled Function to check if the model has audit logging enabled
func isAuditLogEnabled(model interface{}) bool {
	modelType := reflect.TypeOf(model)

	// If the modelType is a pointer, get the underlying type
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Ensure it's a struct
	if modelType.Kind() != reflect.Struct {
		return false
	}

	// Get the type's name
	typeName := modelType.Name()

	// Get the package path using the types package
	pkgPath := modelType.PkgPath()
	return audit.LogModels[pkgPath+typeName]
}

// hasPreSaveHook Check if the model has a PreSave method (custom user hook)
func hasPreSaveHook(model interface{}) bool {
	_, ok := reflect.TypeOf(model).MethodByName("PreSave")
	return ok
}

// hasPostSaveHook Check if the model has a PostSave method (custom user hook)
func hasPostSaveHook(model interface{}) bool {
	_, ok := reflect.TypeOf(model).MethodByName("PostSave")
	return ok
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)

	// Check if the input is a pointer and get the element
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure the value is a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		// Get the JSON and BSON tags
		jsonTag := field.Tag.Get("json")
		bsonTag := field.Tag.Get("bson")

		// Check if the field should be omitted
		if strings.Contains(bsonTag, "omitempty") {
			if isOmitEmpty(value) {
				continue // Skip adding this field if it's empty and has "omitempty"
			}
		}

		// Use JSON tag as the key if present; otherwise, fallback to BSON tag
		if bsonTag != "" {
			splitTag := strings.Split(bsonTag, ",")
			if splitTag[0] == "_id" {
				objectID := value.Interface().(primitive.ObjectID)
				result[splitTag[0]] = objectID.Hex()
				continue
			}
			result[splitTag[0]] = value.Interface()
		} else if jsonTag != "" {
			splitTag := strings.Split(jsonTag, ",")
			result[splitTag[0]] = value.Interface()
		} else {
			result[convertToSnakeCase(field.Name)] = value.Interface()
		}
	}

	return result, nil
}

// isOmitEmpty checks if a value is considered "empty" according to the omitempty rule
func isOmitEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil() // Nil pointer is considered empty
	case reflect.Slice, reflect.Array:
		return value.Len() == 0 // Empty slice/array is considered empty
	case reflect.Map:
		return value.IsNil() || value.Len() == 0 // Nil or empty map is considered empty
	default:
		// For all other types, zero value is considered empty
		zero := reflect.Zero(value.Type())
		return value.Interface() == zero.Interface()
	}
}

func convertToSnakeCase(str string) string {
	var snakeCase string
	runes := []rune(str)
	length := len(runes)

	for i := 0; i < length; i++ {
		// Check if the current character is an uppercase letter
		if unicode.IsUpper(runes[i]) {
			// Add underscore before the uppercase letter unless it's the first character or the previous character is already an underscore
			if i > 0 && runes[i-1] != '_' && (unicode.IsLower(runes[i-1]) || (i+1 < length && unicode.IsLower(runes[i+1]))) {
				snakeCase += "_"
			}
			// Convert the uppercase letter to lowercase
			snakeCase += string(unicode.ToLower(runes[i]))
		} else {
			// Just append the lowercase letters and underscores
			snakeCase += string(runes[i])
		}
	}
	return snakeCase
}

func compareDocumentStates(oldDoc, newDoc map[string]interface{}) map[string]entities.AuditChange {
	changes := make(map[string]entities.AuditChange)

	// Check for keys that are in newDoc (added/modified keys)
	for key, newVal := range newDoc {
		if key == "_id" {
			continue
		}
		oldVal, exists := oldDoc[key]
		if !exists || oldVal != newVal {
			changes[key] = entities.AuditChange{
				Old: fmt.Sprintf("%v", oldVal),
				New: fmt.Sprintf("%v", newVal),
			}
		}
	}

	// Check for keys that are in oldDoc but not in newDoc (deleted keys)
	for key, oldVal := range oldDoc {
		if key == "_id" {
			continue
		}
		if _, exists := newDoc[key]; !exists {
			changes[key] = entities.AuditChange{
				Old: fmt.Sprintf("%v", oldVal),
				New: "", // Key was deleted, so no new value
			}
		}
	}

	return changes
}
