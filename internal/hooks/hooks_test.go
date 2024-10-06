package hooks

import (
	"context"
	audit "github.com/its-own/gaudit/internal/audit_log"
	"github.com/its-own/gaudit/internal/entities"
	in "github.com/its-own/gaudit/pkg"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

// Sample struct for testing
type TestStruct struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Age      int                `json:"age,omitempty" bson:"age,omitempty"`
	IsActive bool               `bson:"is_active"`
	Address  string             `bson:"address,omitempty"`
	EmptyVal string             `bson:"empty_val,omitempty"`
	Unmapped string
}

// TestStruct with PreSave and PostSave methods
type TestHookStruct struct{}

func (t *TestHookStruct) PreSave()  {}
func (t *TestHookStruct) PostSave() {}

// TestStructWithoutHooks without PreSave and PostSave methods
type TestStructWithoutHooks struct{}

// Mock models for testing
type TestModelWithAudit struct {
	in.Inject
}
type TestModelWithoutAudit struct{}
type NonStructModel int

func TestGetContextValue(t *testing.T) {
	// Define the key to be used in context
	const key = "testKey"

	// Create test cases
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "Key exists with non-empty value",
			ctx:      context.WithValue(context.Background(), key, "value"),
			expected: "value",
		},
		{
			name:     "Key exists with empty value",
			ctx:      context.WithValue(context.Background(), key, ""),
			expected: "default",
		},
		{
			name:     "Key does not exist",
			ctx:      context.Background(),
			expected: "default",
		},
		{
			name:     "Key exists with non-string value",
			ctx:      context.WithValue(context.Background(), key, 123), // using int instead of string
			expected: "default",
		},
	}

	// Iterate through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getContextValue(tt.ctx, key)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsAuditLogEnabled(t *testing.T) {
	tests := []struct {
		name     string
		model    interface{}
		expected bool
	}{
		{"TestModelWithAudit", &TestModelWithAudit{}, true},
		{"TestModelWithoutAudit", &TestModelWithoutAudit{}, false},
		{"NonStructModel", NonStructModel(1), false}, // Non-struct input
		{"EmptyStruct", struct{}{}, false},           // Empty struct
	}
	audit.LogModels["github.com/its-own/gaudit/internal/hooksTestModelWithAudit"] = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAuditLogEnabled(tt.model)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func Test_hasPreSaveHook(t *testing.T) {
	tests := []struct {
		name    string
		model   interface{}
		hasHook bool
	}{
		{"With PreSave", &TestHookStruct{}, true},
		{"Without PreSave", &TestStructWithoutHooks{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPreSaveHook(tt.model)
			if got != tt.hasHook {
				t.Errorf("hasPreSaveHook() = %v, want %v", got, tt.hasHook)
			}
		})
	}
}

func Test_hasPostSaveHook(t *testing.T) {
	tests := []struct {
		name    string
		model   interface{}
		hasHook bool
	}{
		{"With PostSave", &TestHookStruct{}, true},
		{"Without PostSave", &TestStructWithoutHooks{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPostSaveHook(tt.model)
			if got != tt.hasHook {
				t.Errorf("hasPostSaveHook() = %v, want %v", got, tt.hasHook)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	// Create test object
	obj := TestStruct{
		ID:       primitive.NewObjectID(),
		Name:     "John Doe",
		Age:      0, // should be omitted due to "omitempty"
		IsActive: true,
		Address:  "123 Main St",
		EmptyVal: "",
		Unmapped: "UnmappedField",
	}

	// Call structToMap to convert to a map
	result, err := structToMap(obj)
	if err != nil {
		t.Fatalf("Error converting struct to map: %v", err)
	}

	// Expected map output
	expected := map[string]interface{}{
		"_id":       obj.ID.Hex(),
		"name":      "John Doe",
		"is_active": true,
		"address":   "123 Main St",
		"unmapped":  "UnmappedField",
	}

	// Check that all expected keys exist in the result
	for key, expectedValue := range expected {
		value, exists := result[key]
		if !exists {
			t.Errorf("Key %s not found in result map", key)
		} else if value != expectedValue {
			t.Errorf("Value for key %s mismatch. Expected %v, got %v", key, expectedValue, value)
		}
	}

	// Ensure omitted fields are not present
	omittedKeys := []string{"age", "empty_val"}
	for _, key := range omittedKeys {
		if _, exists := result[key]; exists {
			t.Errorf("Key %s should be omitted but was found in result", key)
		}
	}
}

func Test_compareDocumentStates(t *testing.T) {
	tests := []struct {
		name         string
		oldDoc       map[string]interface{}
		newDoc       map[string]interface{}
		expectedDiff map[string]entities.AuditChange
	}{
		{
			name: "No changes",
			oldDoc: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			newDoc: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			expectedDiff: map[string]entities.AuditChange{},
		},
		{
			name: "New key added",
			oldDoc: map[string]interface{}{
				"name": "test",
			},
			newDoc: map[string]interface{}{
				"name":  "test",
				"age":   30,
				"email": "test@example.com",
			},
			expectedDiff: map[string]entities.AuditChange{
				"age":   {Old: "<nil>", New: "30"},
				"email": {Old: "<nil>", New: "test@example.com"},
			},
		},
		{
			name: "Value changed",
			oldDoc: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			newDoc: map[string]interface{}{
				"name": "test",
				"age":  35,
			},
			expectedDiff: map[string]entities.AuditChange{
				"age": {Old: "30", New: "35"},
			},
		},
		{
			name: "Key removed in new doc",
			oldDoc: map[string]interface{}{
				"name":  "test",
				"email": "test@example.com",
			},
			newDoc: map[string]interface{}{
				"name": "test",
			},
			expectedDiff: map[string]entities.AuditChange{
				"email": {Old: "test@example.com", New: ""},
			},
		},
		{
			name: "Ignore _id field",
			oldDoc: map[string]interface{}{
				"name": "test",
				"_id":  "some_id",
			},
			newDoc: map[string]interface{}{
				"name": "test",
				"_id":  "new_id",
			},
			expectedDiff: map[string]entities.AuditChange{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := compareDocumentStates(tt.oldDoc, tt.newDoc)
			assert.Equal(t, tt.expectedDiff, diff)
		})
	}
}

func Test_convertToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test cases for regular camel case
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"snake_case", "snake_case"},

		// Test cases for acronyms
		{"userID", "user_id"},
		{"HTTPServer", "http_server"},
		{"MyURL", "my_url"},

		// Test cases for mixed acronyms and camel case
		{"getUserID", "get_user_id"},
		{"getHTTPResponse", "get_http_response"},

		// Edge cases
		{"", ""},
		{"single", "single"},
		{"Already_Snake_Case", "already_snake_case"}, // Shouldn't change
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, but got %s", tt.expected, result)
			}
		})
	}
}

func Test_isOmitEmpty(t *testing.T) {
	// Pointer types
	var ptr *int
	if !isOmitEmpty(reflect.ValueOf(ptr)) {
		t.Errorf("Expected true for nil pointer")
	}

	ptr = new(int)
	if isOmitEmpty(reflect.ValueOf(ptr)) {
		t.Errorf("Expected false for non-nil pointer")
	}

	// Slice types
	var slice []int
	if !isOmitEmpty(reflect.ValueOf(slice)) {
		t.Errorf("Expected true for nil slice")
	}

	slice = []int{}
	if !isOmitEmpty(reflect.ValueOf(slice)) {
		t.Errorf("Expected true for empty slice")
	}

	slice = []int{1, 2, 3}
	if isOmitEmpty(reflect.ValueOf(slice)) {
		t.Errorf("Expected false for non-empty slice")
	}

	// Array types
	var array [0]int
	if !isOmitEmpty(reflect.ValueOf(array)) {
		t.Errorf("Expected true for empty array")
	}

	var newArray = [3]int{1, 2, 3}
	if isOmitEmpty(reflect.ValueOf(newArray)) {
		t.Errorf("Expected false for non-empty array")
	}

	// Map types
	var m map[string]int
	if !isOmitEmpty(reflect.ValueOf(m)) {
		t.Errorf("Expected true for nil map")
	}

	m = make(map[string]int)
	if !isOmitEmpty(reflect.ValueOf(m)) {
		t.Errorf("Expected true for empty map")
	}

	m["key"] = 1
	if isOmitEmpty(reflect.ValueOf(m)) {
		t.Errorf("Expected false for non-empty map")
	}

	// Integer types (default type case)
	var i int
	if !isOmitEmpty(reflect.ValueOf(i)) {
		t.Errorf("Expected true for zero integer")
	}

	i = 10
	if isOmitEmpty(reflect.ValueOf(i)) {
		t.Errorf("Expected false for non-zero integer")
	}

	// String types (default type case)
	var s string
	if !isOmitEmpty(reflect.ValueOf(s)) {
		t.Errorf("Expected true for empty string")
	}

	s = "hello"
	if isOmitEmpty(reflect.ValueOf(s)) {
		t.Errorf("Expected false for non-empty string")
	}
}
