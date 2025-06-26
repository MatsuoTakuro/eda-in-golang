package jetstream

import (
	"context"
	"errors"
	sync "sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eda-in-golang/internal/am"
)

const maxRetries = 5

type stream struct {
	module      string
	streamName  string
	js          jetstream.JetStream
	mu          sync.Mutex
	consumeCtxs []jetstream.ConsumeContext
	logger      zerolog.Logger
}

var _ am.MessageStream = (*stream)(nil)

func NewStream(
	module,
	streamName string,
	js jetstream.JetStream,
	logger zerolog.Logger,
) *stream {
	return &stream{
		module:     module,
		streamName: streamName,
		js:         js,
		logger:     logger,
	}
}

func (s *stream) Publish(ctx context.Context, _ string, rawMsg am.Message) (err error) {
	var data []byte

	metadata, err := structpb.NewStruct(rawMsg.Metadata())
	if err != nil {
		return err
	}

	data, err = proto.Marshal(&StreamMessage{
		Id:       rawMsg.ID(),
		Name:     rawMsg.MessageName(),
		Data:     rawMsg.Data(),
		Metadata: metadata,
		SentAt:   timestamppb.New(rawMsg.SentAt()),
	})
	if err != nil {
		return
	}

	var p jetstream.PubAckFuture
	p, err = s.js.PublishMsgAsync(&nats.Msg{
		Subject: rawMsg.Subject(),
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
				var logEvt *zerolog.Event
				if rawMsg.MessageName() == am.FailureReply {
					logEvt = s.logger.Error()
				} else {
					logEvt = s.logger.Info()
				}
				logEvt.
					Any(moduleField, s.module).
					Any(subjectField, rawMsg.Subject()).
					Any(msgNameField, rawMsg.MessageName()).
					Any(msgIDField, rawMsg.ID()).
					Msg("acked publishing message")
				return
			case <-future.Err(): // error ignored; try again
				// TODO add some variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					// TODO do more than give up
					s.logger.Error().
						Any(moduleField, s.module).
						Any(subjectField, rawMsg.Subject()).
						Any(msgNameField, rawMsg.MessageName()).
						Any(msgIDField, rawMsg.ID()).
						Err(err).
						Msg("gave up publishing message")
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					// TODO do more than give up
					s.logger.Error().
						Any(moduleField, s.module).
						Any(subjectField, rawMsg.Subject()).
						Any(msgNameField, rawMsg.MessageName()).
						Any(msgIDField, rawMsg.ID()).
						Err(err).
						Msg("failed to publish message")
					return
				}
			}
		}
	}(p, maxRetries)

	return nil
}

var doNothing = func() error { return nil }

func (s *stream) Subscribe(topicName string, handler am.MessageHandler, options ...am.SubscriberOption) (unsubscribe func() error, err error) {

	s.mu.Lock()
	defer s.mu.Unlock()

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
		return doNothing, err
	}

	c, err := sm.CreateOrUpdateConsumer(ctx, cfg)
	if err != nil {
		return doNothing, err
	}
	// Consume receives messages from the stream and calls the handler for each one.
	// JetStream delivers messages sequentially to this callback,
	// so parallelism must be implemented explicitly using goroutines or a worker pool,
	// such as wrapping the handler in a goroutine or using Messages() to dispatch with workers.
	cc, err := c.Consume(func(msg jetstream.Msg) {
		s.handleMsg(subCfg, handler)(msg)
	})
	if err != nil {
		return doNothing, err
	}
	s.consumeCtxs = append(s.consumeCtxs, cc)

	unsubscribe = func() error {
		cc.Drain()
		return nil
	}

	return unsubscribe, nil
}

func (s *stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler) func(jetstream.Msg) {
	return func(natsMsg jetstream.Msg) {
		var err error

		m := &StreamMessage{}
		err = proto.Unmarshal(natsMsg.Data(), m)
		if err != nil {
			// TODO Nak? ... logging?
			s.logger.Error().
				Any(moduleField, s.module).
				Err(err).
				Msg("failed to unmarshal received message")
			return
		}

		msg := &rawMessage{
			id:       m.GetId(),
			name:     m.GetName(),
			subject:  natsMsg.Subject(),
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
				s.logger.Error().
					Any(moduleField, s.module).
					Any(subjectField, msg.Subject()).
					Any(msgNameField, msg.MessageName()).
					Any(msgIDField, msg.ID()).
					Err(err).
					Msg("failed to auto-ack received message")
			}
		}

		select {
		case err = <-errc:
			s.handleMsgResult(msg, err)

		case <-wCtx.Done():
			s.logger.Error().
				Any(moduleField, s.module).
				Any(subjectField, msg.Subject()).
				Any(msgNameField, msg.MessageName()).
				Any(msgIDField, msg.ID()).
				Err(wCtx.Err()).
				Msg("timeout for handling received message")
		}
	}
}

func (s *stream) handleMsgResult(msg *rawMessage, err error) {
	switch {
	case err == nil:
		if ackErr := msg.Ack(); ackErr != nil {
			s.logger.Error().
				Any(moduleField, s.module).
				Any(subjectField, msg.Subject()).
				Any(msgNameField, msg.MessageName()).
				Any(msgIDField, msg.ID()).
				Err(ackErr).
				Msg("failed to ack received message")
			return
		}
		var logEvt *zerolog.Event
		if msg.MessageName() == am.FailureReply {
			logEvt = s.logger.Error()
		} else {
			logEvt = s.logger.Info()
		}
		logEvt.
			Any(moduleField, s.module).
			Any(subjectField, msg.Subject()).
			Any(msgNameField, msg.MessageName()).
			Any(msgIDField, msg.ID()).
			Msg("acked received message")
		return

	case errors.Is(err, am.ErrMessageSkipped):
		if ackErr := msg.Ack(); ackErr != nil {
			s.logger.Error().
				Any(moduleField, s.module).
				Any(subjectField, msg.Subject()).
				Any(msgNameField, msg.MessageName()).
				Any(msgIDField, msg.ID()).
				Err(ackErr).
				Msg("failed to ack received message for skipping")
			return
		}
		s.logger.Info().
			Any(moduleField, s.module).
			Any(subjectField, msg.Subject()).
			Any(msgNameField, msg.MessageName()).
			Any(msgIDField, msg.ID()).
			Msg("skipped handling received message")
		return

	default:
		if nakErr := msg.NAck(); nakErr != nil {
			err = errors.Join(err, nakErr)
			s.logger.Error().
				Any(moduleField, s.module).
				Any(subjectField, msg.Subject()).
				Any(msgNameField, msg.MessageName()).
				Any(msgIDField, msg.ID()).
				Err(err).
				Msg("failed to nack received message")
			return
		}
		s.logger.Error().
			Any(moduleField, s.module).
			Any(subjectField, msg.Subject()).
			Any(msgNameField, msg.MessageName()).
			Any(msgIDField, msg.ID()).
			Err(err).
			Msg("nacked received message")
		return
	}
}

const (
	moduleField  = "mod"
	subjectField = "sub"
	msgNameField = "name"
	msgIDField   = "id"
)

func (s *stream) UnsubscribeAll() error {
	for _, ctx := range s.consumeCtxs {
		ctx.Drain()
	}

	return nil
}
