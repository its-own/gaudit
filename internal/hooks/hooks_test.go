package hooks

import (
	"testing"
)

func Test_Demo(t *testing.T) {
	//tests := []struct {
	//	name         string
	//	oldDoc       map[string]interface{}
	//	newDoc       map[string]interface{}
	//	expectedDiff map[string]entities.AuditChange
	//}{
	//	{
	//		name: "No changes",
	//		oldDoc: map[string]interface{}{
	//			"name": "test",
	//			"age":  30,
	//		},
	//		newDoc: map[string]interface{}{
	//			"name": "test",
	//			"age":  30,
	//		},
	//		expectedDiff: map[string]entities.AuditChange{},
	//	},
	//	{
	//		name: "New key added",
	//		oldDoc: map[string]interface{}{
	//			"name": "test",
	//		},
	//		newDoc: map[string]interface{}{
	//			"name":  "test",
	//			"age":   30,
	//			"email": "test@example.com",
	//		},
	//		expectedDiff: map[string]entities.AuditChange{
	//			"age":   {Old: "<nil>", New: "30"},
	//			"email": {Old: "<nil>", New: "test@example.com"},
	//		},
	//	},
	//	{
	//		name: "Value changed",
	//		oldDoc: map[string]interface{}{
	//			"name": "test",
	//			"age":  30,
	//		},
	//		newDoc: map[string]interface{}{
	//			"name": "test",
	//			"age":  35,
	//		},
	//		expectedDiff: map[string]entities.AuditChange{
	//			"age": {Old: "30", New: "35"},
	//		},
	//	},
	//	{
	//		name: "Key removed in new doc",
	//		oldDoc: map[string]interface{}{
	//			"name":  "test",
	//			"email": "test@example.com",
	//		},
	//		newDoc: map[string]interface{}{
	//			"name": "test",
	//		},
	//		expectedDiff: map[string]entities.AuditChange{
	//			"email": {Old: "test@example.com", New: "<nil>"},
	//		},
	//	},
	//	{
	//		name: "Ignore _id field",
	//		oldDoc: map[string]interface{}{
	//			"name": "test",
	//			"_id":  "some_id",
	//		},
	//		newDoc: map[string]interface{}{
	//			"name": "test",
	//			"_id":  "new_id",
	//		},
	//		expectedDiff: map[string]entities.AuditChange{},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		diff := compareDocumentStates(tt.oldDoc, tt.newDoc)
	//		assert.Equal(t, tt.expectedDiff, diff)
	//	})
	//}
}
