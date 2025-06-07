package sec

type StepOption[T any] func(step *step[T])

func WithNormal[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.actions[normal] = fn
	}
}

func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.actions[compensating] = fn
	}
}

func OnNormalReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.replyHandlers[normal][replyName] = fn
	}
}

func OnCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.replyHandlers[compensating][replyName] = fn
	}
}
