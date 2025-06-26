package am

type MessagePublisherMiddleware = func(next MessagePublisher) MessagePublisher

func messagePublisherWithMiddleware(
	msgPublisher MessagePublisher,
	mws ...MessagePublisherMiddleware,
) MessagePublisher {
	return applyMiddleware(msgPublisher, mws...)
}

type MessageHandlerMiddleware = func(next MessageHandler) MessageHandler

func messageHandlerWithMiddleware(
	msgHandler MessageHandler,
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return applyMiddleware(msgHandler, mws...)
}

func applyMiddleware[T any, M func(T) T](target T, mws ...M) T {
	h := target
	// apply the middlewares in reverse order
	// so that the first middleware is the last to apply
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
