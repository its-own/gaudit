package in

import "context"

// Hook interface for custom hooks
type Hook interface {
	PreSave(ctx context.Context, model interface{}, filter interface{}, col, ops, docId string)
	PostSave(ctx context.Context, model interface{}, filter interface{}, col, ops, docId string)
}

type Inject struct {
}
