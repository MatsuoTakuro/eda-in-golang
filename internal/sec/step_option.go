package sec

type StepOption[T any] func(step *step[T])

func WithAction[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.actions[notCompensating] = fn
	}
}

func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.actions[isCompensating] = fn
	}
}

func OnActionReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.handlers[notCompensating][replyName] = fn
	}
}

func OnCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *step[T]) {
		step.handlers[isCompensating][replyName] = fn
	}
}
