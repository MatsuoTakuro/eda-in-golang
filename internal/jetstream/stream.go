package jetstream

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/am"
)

const maxRetries = 5

type stream struct {
	streamName  string
	js          jetstream.JetStream
	consumeCtxs []jetstream.ConsumeContext
}

var _ am.MessageStream[am.RawMessage, am.RawMessage] = (*stream)(nil)

func NewStream(streamName string, js jetstream.JetStream) *stream {
	return &stream{
		streamName: streamName,
		js:         js,
	}
}

func (s *stream) Publish(ctx context.Context, topicName string, rawMsg am.RawMessage) (err error) {
	var data []byte

	data, err = proto.Marshal(&StreamMessage{
		Id:   rawMsg.ID(),
		Name: rawMsg.MessageName(),
		Data: rawMsg.Data(),
	})
	if err != nil {
		return err
	}

	var p jetstream.PubAckFuture
	p, err = s.js.PublishMsgAsync(&nats.Msg{
		Subject: topicName,
		Data:    data,
	}, jetstream.WithMsgID(rawMsg.ID()))
	if err != nil {
		return err
	}

	// retry a handful of times to publish the messages
	go func(future jetstream.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok(): // publish acknowledged
				log.Printf("acknowledged message published: (%s: %s)", rawMsg.MessageName(), rawMsg.ID())
				return
			case <-future.Err(): // error ignored; try again
				// TODO add some variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					// TODO do more than give up
					log.Printf("gave up publishing message after %d retries: (%s: %s): %v", maxRetries, rawMsg.MessageName(), rawMsg.ID(), err)
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					// TODO do more than give up
					log.Printf("failed to publish message: (%s: %s): %v", rawMsg.MessageName(), rawMsg.ID(), err)
					return
				}
			}
		}
	}(p, maxRetries)

	return nil
}

func (s *stream) Subscribe(topicName string, handler am.MessageHandler[am.RawMessage], options ...am.SubscriberOption) error {
	var err error

	subCfg := am.NewSubscriberConfig(options)

	cfg := jetstream.ConsumerConfig{
		MaxDeliver:    subCfg.MaxRedeliver(),
		FilterSubject: topicName,
	}

	// Durable sets the consumer name to make it persistent across service restarts.
	// If not set, the consumer will be ephemeral and deleted automatically.
	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.Durable = groupName
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		cfg.AckPolicy = jetstream.AckExplicitPolicy
		cfg.AckWait = subCfg.AckWait()
	} else {
		cfg.AckPolicy = jetstream.AckNonePolicy
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sm, err := s.js.Stream(ctx, s.streamName)
	if err != nil {
		return err
	}

	c, err := sm.CreateOrUpdateConsumer(ctx, cfg)
	if err != nil {
		return err
	}
	// Consume receives messages from the stream and calls the handler for each one.
	// JetStream delivers messages sequentially to this callback,
	// so parallelism must be implemented explicitly using goroutines or a worker pool,
	// such as wrapping the handler in a goroutine or using Messages() to dispatch with workers.
	cc, err := c.Consume(func(msg jetstream.Msg) {
		s.handleMsg(subCfg, handler)(msg)
	})
	if err != nil {
		return err
	}
	s.consumeCtxs = append(s.consumeCtxs, cc)

	return nil
}

func (s *stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler[am.RawMessage]) func(jetstream.Msg) {
	return func(natsMsg jetstream.Msg) {
		var err error

		m := &StreamMessage{}
		err = proto.Unmarshal(natsMsg.Data(), m)
		if err != nil {
			// TODO Nak? ... logging?
			log.Printf("failed to unmarshal message received: %v", err)
			return
		}

		msg := &rawMessage{
			id:       m.GetId(),
			name:     m.GetName(),
			data:     m.GetData(),
			acked:    false,
			ackFn:    func() error { return natsMsg.Ack() },
			nackFn:   func() error { return natsMsg.Nak() },
			extendFn: func() error { return natsMsg.InProgress() },
			killFn:   func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				// TODO logging?
				log.Printf("failed to auto-ack message received: %v", err)
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					// TODO logging?
					log.Printf("failed to ack message received (%s: %s): %v", msg.MessageName(), msg.ID(), ackErr)
				}
				log.Printf("acknowledged message received: (%s: %s)", msg.MessageName(), msg.ID())
				return
			}
			if nakErr := msg.NAck(); nakErr != nil {
				// TODO logging?
				log.Printf("failed to nack message received (%s: %s): %v", msg.MessageName(), msg.ID(), nakErr)

			}
		case <-wCtx.Done():
			// TODO logging?
			log.Printf("timeout for handling message received (%s: %s): %v", msg.MessageName(), msg.ID(), wCtx.Err())
			return
		}
	}
}

// TODO: Not sure if we should call this before nats connection is drained.
func (s *stream) Drain() error {
	for _, ctx := range s.consumeCtxs {
		ctx.Drain()
	}

	return nil
}
