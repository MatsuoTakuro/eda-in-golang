package sec

import "context"

type Store interface {
	Load(ctx context.Context, sagaName, sagaID string) (*Context[[]byte], error)
	Save(ctx context.Context, sagaName string, sagaCtx *Context[[]byte]) error
}
