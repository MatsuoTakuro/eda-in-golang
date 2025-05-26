package jetstream

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/am"
)

const maxRetries = 5

type stream struct {
	module      string
	streamName  string
	js          jetstream.JetStream
	consumeCtxs []jetstream.ConsumeContext
}

var _ am.MessageStream[am.RawMessage, am.RawMessage] = (*stream)(nil)

func NewStream(module, streamName string, js jetstream.JetStream) *stream {
	return &stream{
		module:     module,
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
				log.Printf("%s acked publishing message: (%s: %s)", s.module, rawMsg.MessageName(), rawMsg.ID())
				return
			case <-future.Err(): // error ignored; try again
				// TODO add some variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					// TODO do more than give up
					log.Printf("%s gave up publishing message after %d retries: (%s: %s): %v", s.module, maxRetries, rawMsg.MessageName(), rawMsg.ID(), err)
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					// TODO do more than give up
					log.Printf("%s failed to publish message: (%s: %s): %v", s.module, rawMsg.MessageName(), rawMsg.ID(), err)
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
			log.Printf("failed to unmarshal received message: %v", err)
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

		errc := make(chan error) // no buffer, so we can wait for the handler to finish
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				// TODO logging?
				log.Printf("%s failed to auto-ack received message: %v", s.module, err)
			}
		}

		select {
		case err = <-errc:
			handleMsgResult(s.module, msg, err)

		case <-wCtx.Done():
			log.Printf("%s timeout for handling received message (%s: %s): %v", s.module, msg.MessageName(), msg.ID(), wCtx.Err())
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

func handleMsgResult(module string, msg *rawMessage, err error) {
	switch {
	case err == nil:
		if ackErr := msg.Ack(); ackErr != nil {
			log.Printf("%s failed to ack received message (%s: %s): %v", module, msg.MessageName(), msg.ID(), ackErr)
			return
		}
		log.Printf("%s acked received message: (%s: %s)", module, msg.MessageName(), msg.ID())
		return

	case errors.Is(err, am.ErrMessageSkipped):
		if ackErr := msg.Ack(); ackErr != nil {
			log.Printf("%s failed to ack received message (%s: %s): %v", module, msg.MessageName(), msg.ID(), ackErr)
			return
		}
		log.Printf("%s skipped handling received message (%s: %s)", module, msg.MessageName(), msg.ID())
		return

	default:
		if nakErr := msg.NAck(); nakErr != nil {
			err = errors.Join(err, nakErr)
			log.Printf("%s failed to nack received message (%s: %s): %v", module, msg.MessageName(), msg.ID(), err)
			return
		}
		log.Printf("%s nacked received message: (%s: %s): %v", module, msg.MessageName(), msg.ID(), err)
		return
	}
}
