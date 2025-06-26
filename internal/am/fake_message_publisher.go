package am

import (
	"context"

	"github.com/stackus/errors"
)

type fakeMessage struct {
	subject string
	payload Message
}

type FakeMessagePublisher struct {
	messages []fakeMessage
}

var _ MessagePublisher = (*FakeMessagePublisher)(nil)

func NewFakeMessagePublisher() *FakeMessagePublisher {
	return &FakeMessagePublisher{
		messages: []fakeMessage{},
	}
}

func (p *FakeMessagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	p.messages = append(p.messages, fakeMessage{topicName, msg})
	return nil
}

func (p *FakeMessagePublisher) Reset() {
	p.messages = []fakeMessage{}
}

func (p *FakeMessagePublisher) Last() (string, Message, error) {
	var msg Message
	if len(p.messages) == 0 {
		return "", msg, errors.ErrNotFound.Msg("no messages have been published")
	}

	last := p.messages[len(p.messages)-1]
	return last.subject, last.payload, nil
}
