package sec

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type Orchestrator[T any] interface {
	// Start begins the Saga execution from the first step.
	Start(ctx context.Context, id string, data T) error
	// ReplyTopic returns the topic on which this orchestrator listens for replies.
	ReplyTopic() string
	// HandleReply processes a reply message and determines the next Saga step.
	HandleReply(ctx context.Context, reply ddd.Reply) error
}

type orchestrator[T any] struct {
	saga      Saga[T]
	repo      repository[T]
	publisher am.CommandPublisher
}

var _ Orchestrator[any] = (*orchestrator[any])(nil)

func NewOrchestrator[T any](
	saga Saga[T],
	repo repository[T],
	publisher am.CommandPublisher,
) orchestrator[T] {
	return orchestrator[T]{
		saga:      saga,
		repo:      repo,
		publisher: publisher,
	}
}

const noStepsExecuted = -1

func (o orchestrator[T]) Start(ctx context.Context, id string, data T) error {
	sagaCtx := &Context[T]{
		ID:   id,
		Data: data,
		Step: noStepsExecuted,
	}

	err := o.repo.Save(ctx, o.saga.Name(), sagaCtx)
	if err != nil {
		return err
	}

	result := o.execute(ctx, sagaCtx)
	if result.err != nil {
		return err
	}

	return o.processResult(ctx, result)
}

func (o orchestrator[T]) ReplyTopic() string {
	return o.saga.ReplyTopic()
}

func (o orchestrator[T]) HandleReply(ctx context.Context, reply ddd.Reply) error {
	sagaID, sagaName := o.getSagaInfoFromReply(reply)
	if sagaID == "" || sagaName == "" || sagaName != o.saga.Name() {
		// returning nil to drop bad replies
		return nil
	}

	sagaCtx, err := o.repo.Load(ctx, o.saga.Name(), sagaID)
	if err != nil {
		return err
	}

	result, err := o.handleReply(ctx, sagaCtx, reply)
	if err != nil {
		return err
	}

	return o.processResult(ctx, result)
}

func (o orchestrator[T]) handleReply(ctx context.Context, sagaCtx *Context[T], reply ddd.Reply) (stepResult[T], error) {
	// Get the current step based on the saga context's step index.
	step := o.saga.getSteps()[sagaCtx.Step]

	err := step.handleReply(ctx, sagaCtx, reply)
	if err != nil {
		return stepResult[T]{}, err
	}

	var isSuccessful bool
	if outcome, ok := reply.Metadata().Get(am.ReplyOutcomeHdr).(string); !ok {
		isSuccessful = false
	} else {
		isSuccessful = (outcome == am.OutcomeSuccess)
	}

	if isSuccessful {
		return o.execute(ctx, sagaCtx), nil
	}

	if sagaCtx.IsCompensating {
		return stepResult[T]{}, errors.ErrInternal.Msg("received failed reply but already compensating")
	}

	sagaCtx.compensate()
	return o.execute(ctx, sagaCtx), nil
}

func (o orchestrator[T]) execute(ctx context.Context, sagaCtx *Context[T]) stepResult[T] {

	var direction = 1 // forward direction for normal steps
	if sagaCtx.IsCompensating {
		direction = -1 // backward direction for compensating steps
	}

	steps := o.saga.getSteps()
	stepCap := len(steps)

	var next Step[T]
	// delta tracks how many steps to move from current position to the next invocable step
	// whether moving forward or backward.
	var delta = 1
	for i := sagaCtx.Step + direction;  // start from the next step
	noStepsExecuted < i && i < stepCap; // ensure we stay within bounds
	i += direction {                    // increment by direction
		next = steps[i]
		if next != nil && next.isInvocable(sagaCtx.IsCompensating) {
			break
		}
		delta += 1 // track the number of steps moved (forward or backward)
	}

	if next == nil {
		sagaCtx.complete()
		return stepResult[T]{ctx: sagaCtx}
	}

	sagaCtx.advance(delta)

	return next.executeAction(ctx, sagaCtx)
}

func (o orchestrator[T]) processResult(ctx context.Context, result stepResult[T]) (err error) {
	if result.next != nil {
		err = o.publishCommand(ctx, result)
		if err != nil {
			return
		}
	}

	return o.repo.Save(ctx, o.saga.Name(), result.ctx)
}

func (o orchestrator[T]) publishCommand(ctx context.Context, result stepResult[T]) error {
	cmd := result.next

	cmd.Metadata().Set(am.CommandReplyChannelHdr, o.saga.ReplyTopic())
	cmd.Metadata().Set(SagaCommandIDHdr, result.ctx.ID)
	cmd.Metadata().Set(SagaCommandNameHdr, o.saga.Name())

	return o.publisher.Publish(ctx, cmd.Destination(), cmd)
}

func (o orchestrator[T]) getSagaInfoFromReply(reply ddd.Reply) (string, string) {
	var ok bool
	var sagaID, sagaName string

	if sagaID, ok = reply.Metadata().Get(SagaReplyIDHdr).(string); !ok {
		return "", ""
	}

	if sagaName, ok = reply.Metadata().Get(SagaReplyNameHdr).(string); !ok {
		return "", ""
	}

	return sagaID, sagaName
}
