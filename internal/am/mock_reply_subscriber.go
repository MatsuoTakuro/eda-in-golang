// Code generated by mockery v2.53.4. DO NOT EDIT.

package am

import mock "github.com/stretchr/testify/mock"

// MockReplySubscriber is an autogenerated mock type for the ReplySubscriber type
type MockReplySubscriber[A AckableMessage] struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: topicName, handler, options
func (_m *MockReplySubscriber[A]) Subscribe(topicName string, handler MessageHandler[ReplyMessage], options ...SubscriberOption) error {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, topicName, handler)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Subscribe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, MessageHandler[ReplyMessage], ...SubscriberOption) error); ok {
		r0 = rf(topicName, handler, options...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with no fields
func (_m *MockReplySubscriber[A]) Unsubscribe() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Unsubscribe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockReplySubscriber creates a new instance of MockReplySubscriber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockReplySubscriber[A AckableMessage](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockReplySubscriber[A] {
	mock := &MockReplySubscriber[A]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
