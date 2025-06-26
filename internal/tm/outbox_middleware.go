package tm

import (
	"context"

	"errors"

	"eda-in-golang/internal/am"
)

type OutboxStore interface {
	Save(ctx context.Context, msg am.Message) error
	FindUnpublished(ctx context.Context, limit int) ([]am.Message, error)
	MarkPublished(ctx context.Context, ids ...string) error
}

func WithOutboxStream(store OutboxStore) am.RawMessageStreamMiddleware {
	o := outbox{store: store}

	return func(stream am.RawMessageStream) am.RawMessageStream {
		o.RawMessageStream = stream

		return o
	}
}

type outbox struct {
	am.RawMessageStream
	store OutboxStore
}

var _ am.RawMessageStream = (*outbox)(nil)

func (o outbox) Publish(ctx context.Context, _ string, msg am.Message) error {
	err := o.store.Save(ctx, msg)

	var errDupe ErrDuplicateMessage
	if errors.As(err, &errDupe) {
		return nil
	}

	return err
}
