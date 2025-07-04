// Code generated by mockery v2.53.4. DO NOT EDIT.

package am

import (
	context "context"
	ddd "eda-in-golang/internal/ddd"

	mock "github.com/stretchr/testify/mock"
)

// MockCommandPublisher is an autogenerated mock type for the CommandPublisher type
type MockCommandPublisher[T interface{}] struct {
	mock.Mock
}

// Publish provides a mock function with given fields: ctx, topicName, v
func (_m *MockCommandPublisher[T]) Publish(ctx context.Context, topicName string, v ddd.Command) error {
	ret := _m.Called(ctx, topicName, v)

	if len(ret) == 0 {
		panic("no return value specified for Publish")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ddd.Command) error); ok {
		r0 = rf(ctx, topicName, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockCommandPublisher creates a new instance of MockCommandPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommandPublisher[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommandPublisher[T] {
	mock := &MockCommandPublisher[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
