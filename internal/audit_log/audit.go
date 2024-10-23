package audit

import (
	"context"
)

// LogModels Registry for audit-log-enabled models
var LogModels = make(map[string]bool)

func init() {
	err := WatchAndInjectHooks(context.Background())
	if err != nil {
		panic(err)
	}
}

// RegisterModel Register the model for audit logging
func RegisterModel(key string) {
	LogModels[key] = true
}
