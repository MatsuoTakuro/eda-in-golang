// Code generated by mockery v2.53.4. DO NOT EDIT.

package customerspb

import mock "github.com/stretchr/testify/mock"

// MockUnsafeCustomersServiceServer is an autogenerated mock type for the UnsafeCustomersServiceServer type
type MockUnsafeCustomersServiceServer struct {
	mock.Mock
}

// mustEmbedUnimplementedCustomersServiceServer provides a mock function with no fields
func (_m *MockUnsafeCustomersServiceServer) mustEmbedUnimplementedCustomersServiceServer() {
	_m.Called()
}

// NewMockUnsafeCustomersServiceServer creates a new instance of MockUnsafeCustomersServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUnsafeCustomersServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUnsafeCustomersServiceServer {
	mock := &MockUnsafeCustomersServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
