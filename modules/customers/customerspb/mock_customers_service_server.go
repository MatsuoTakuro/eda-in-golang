// Code generated by mockery v2.53.4. DO NOT EDIT.

package customerspb

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCustomersServiceServer is an autogenerated mock type for the CustomersServiceServer type
type MockCustomersServiceServer struct {
	mock.Mock
}

// AuthorizeCustomer provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) AuthorizeCustomer(_a0 context.Context, _a1 *AuthorizeCustomerRequest) (*AuthorizeCustomerResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AuthorizeCustomer")
	}

	var r0 *AuthorizeCustomerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *AuthorizeCustomerRequest) (*AuthorizeCustomerResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *AuthorizeCustomerRequest) *AuthorizeCustomerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AuthorizeCustomerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *AuthorizeCustomerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangeSmsNumber provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) ChangeSmsNumber(_a0 context.Context, _a1 *ChangeSmsNumberRequest) (*ChangeSmsNumberResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ChangeSmsNumber")
	}

	var r0 *ChangeSmsNumberResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ChangeSmsNumberRequest) (*ChangeSmsNumberResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ChangeSmsNumberRequest) *ChangeSmsNumberResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ChangeSmsNumberResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ChangeSmsNumberRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DisableCustomer provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) DisableCustomer(_a0 context.Context, _a1 *DisableCustomerRequest) (*DisableCustomerResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DisableCustomer")
	}

	var r0 *DisableCustomerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *DisableCustomerRequest) (*DisableCustomerResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *DisableCustomerRequest) *DisableCustomerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DisableCustomerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *DisableCustomerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnableCustomer provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) EnableCustomer(_a0 context.Context, _a1 *EnableCustomerRequest) (*EnableCustomerResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for EnableCustomer")
	}

	var r0 *EnableCustomerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *EnableCustomerRequest) (*EnableCustomerResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *EnableCustomerRequest) *EnableCustomerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*EnableCustomerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *EnableCustomerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomer provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) GetCustomer(_a0 context.Context, _a1 *GetCustomerRequest) (*GetCustomerResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetCustomer")
	}

	var r0 *GetCustomerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetCustomerRequest) (*GetCustomerResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetCustomerRequest) *GetCustomerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetCustomerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetCustomerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterCustomer provides a mock function with given fields: _a0, _a1
func (_m *MockCustomersServiceServer) RegisterCustomer(_a0 context.Context, _a1 *RegisterCustomerRequest) (*RegisterCustomerResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RegisterCustomer")
	}

	var r0 *RegisterCustomerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *RegisterCustomerRequest) (*RegisterCustomerResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *RegisterCustomerRequest) *RegisterCustomerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RegisterCustomerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *RegisterCustomerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedCustomersServiceServer provides a mock function with no fields
func (_m *MockCustomersServiceServer) mustEmbedUnimplementedCustomersServiceServer() {
	_m.Called()
}

// NewMockCustomersServiceServer creates a new instance of MockCustomersServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCustomersServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCustomersServiceServer {
	mock := &MockCustomersServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
