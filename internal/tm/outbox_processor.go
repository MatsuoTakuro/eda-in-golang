package tm

import (
	"context"
	"time"

	"eda-in-golang/internal/am"
)

type OutboxProcessor interface {
	Start(ctx context.Context) error
}

type outboxProcessor struct {
	publisher am.RawMessagePublisher
	store     OutboxStore
}

var _ OutboxProcessor = (*outboxProcessor)(nil)

func NewOutboxProcessor(publisher am.RawMessagePublisher, store OutboxStore) outboxProcessor {
	return outboxProcessor{
		publisher: publisher,
		store:     store,
	}
}

func (p outboxProcessor) Start(ctx context.Context) error {
	errC := make(chan error)

	go func() {
		errC <- p.processMessages(ctx)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errC:
		return err
	}
}

const (
	messageLimit    = 50
	pollingInterval = 500 * time.Millisecond
)

func (p outboxProcessor) processMessages(ctx context.Context) error {
	// set up 0 duration timer to poll immediately for first batch
	timer := time.NewTimer(0)
	for {
		msgs, err := p.store.FindUnpublished(ctx, messageLimit)
		if err != nil {
			return err
		}

		if len(msgs) > 0 {
			ids := make([]string, len(msgs))
			for i, msg := range msgs {
				ids[i] = msg.ID()
				err = p.publisher.Publish(ctx, msg.Subject(), msg)
				if err != nil {
					return err
				}
			}
			err = p.store.MarkPublished(ctx, ids...)
			if err != nil {
				return err
			}

			// poll again immediately
			continue
		}

		// if no messages found, stop the timer if it's running
		if !timer.Stop() {
			select {
			case <-timer.C: // drain the channel if it was already running
			default:
				// do nothing, channel is empty
			}
		}

		// wait a short time before polling again
		timer.Reset(pollingInterval)

		select {
		// check if the context is done to exit gracefully
		case <-ctx.Done():
			return nil
		// wait for the reset timer to expire before polling again
		case <-timer.C:
		}
	}
}
