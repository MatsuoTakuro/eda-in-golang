// Code generated by mockery v2.53.4. DO NOT EDIT.

package basketspb

import mock "github.com/stretchr/testify/mock"

// MockUnsafeBasketServiceServer is an autogenerated mock type for the UnsafeBasketServiceServer type
type MockUnsafeBasketServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedBasketServiceServer provides a mock function with no fields
func (_m *MockUnsafeBasketServiceServer) mustEmbedUnimplementedBasketServiceServer() {
	_m.Called()
}

// NewMockUnsafeBasketServiceServer creates a new instance of MockUnsafeBasketServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUnsafeBasketServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUnsafeBasketServiceServer {
	mock := &MockUnsafeBasketServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
