package sec

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/registry"
)

type Repository interface {
	Load(ctx context.Context, sagaName, sagaID string) (*Context[any], error)
	Save(ctx context.Context, sagaName string, sagaCtx *Context[any]) error
}

type repository[T any] struct {
	reg   registry.Registry
	store Store
}

var _ Repository = (*repository[any])(nil)

func NewRepository[T any](reg registry.Registry, store Store) repository[T] {
	return repository[T]{
		reg:   reg,
		store: store,
	}
}

func (r repository[T]) Load(ctx context.Context, sagaName, sagaID string) (*Context[T], error) {
	byteCtx, err := r.store.Load(ctx, sagaName, sagaID)
	if err != nil {
		return nil, err
	}

	v, err := r.reg.Deserialize(sagaName, byteCtx.Data)
	if err != nil {
		return nil, err
	}

	var data T
	var ok bool
	if data, ok = v.(T); !ok {
		return nil, errors.ErrInternal.Msgf("%T is not the expected type %T", v, data)
	}

	return &Context[T]{
		ID:             byteCtx.ID,
		Data:           data,
		Step:           byteCtx.Step,
		Done:           byteCtx.Done,
		IsCompensating: byteCtx.IsCompensating,
	}, nil
}

func (r repository[T]) Save(ctx context.Context, sagaName string, sagaCtx *Context[T]) error {
	data, err := r.reg.Serialize(sagaName, sagaCtx.Data)
	if err != nil {
		return err
	}

	return r.store.Save(ctx, sagaName, &Context[[]byte]{
		ID:             sagaCtx.ID,
		Data:           data,
		Step:           sagaCtx.Step,
		Done:           sagaCtx.Done,
		IsCompensating: sagaCtx.IsCompensating,
	})
}
