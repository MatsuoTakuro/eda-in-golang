package sec

import "context"

// RawStore is used for saga repositories that load or save raw byte slices in saga contexts.
type RawStore interface {
	Load(ctx context.Context, sagaName, sagaID string) (*Context[[]byte], error)
	Save(ctx context.Context, sagaName string, sagaCtx *Context[[]byte]) error
}
