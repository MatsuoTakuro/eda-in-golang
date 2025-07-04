// Code generated by mockery v2.53.4. DO NOT EDIT.

package basketspb

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockBasketServiceServer is an autogenerated mock type for the BasketServiceServer type
type MockBasketServiceServer struct {
	mock.Mock
}

// AddItem provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) AddItem(_a0 context.Context, _a1 *AddItemRequest) (*AddItemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AddItem")
	}

	var r0 *AddItemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *AddItemRequest) (*AddItemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *AddItemRequest) *AddItemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AddItemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *AddItemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CancelBasket provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) CancelBasket(_a0 context.Context, _a1 *CancelBasketRequest) (*CancelBasketResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CancelBasket")
	}

	var r0 *CancelBasketResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CancelBasketRequest) (*CancelBasketResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CancelBasketRequest) *CancelBasketResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CancelBasketResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CancelBasketRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckoutBasket provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) CheckoutBasket(_a0 context.Context, _a1 *CheckoutBasketRequest) (*CheckoutBasketResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CheckoutBasket")
	}

	var r0 *CheckoutBasketResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CheckoutBasketRequest) (*CheckoutBasketResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CheckoutBasketRequest) *CheckoutBasketResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CheckoutBasketResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CheckoutBasketRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBasket provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) GetBasket(_a0 context.Context, _a1 *GetBasketRequest) (*GetBasketResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetBasket")
	}

	var r0 *GetBasketResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetBasketRequest) (*GetBasketResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetBasketRequest) *GetBasketResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetBasketResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetBasketRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveItem provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) RemoveItem(_a0 context.Context, _a1 *RemoveItemRequest) (*RemoveItemResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RemoveItem")
	}

	var r0 *RemoveItemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *RemoveItemRequest) (*RemoveItemResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *RemoveItemRequest) *RemoveItemResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RemoveItemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *RemoveItemRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartBasket provides a mock function with given fields: _a0, _a1
func (_m *MockBasketServiceServer) StartBasket(_a0 context.Context, _a1 *StartBasketRequest) (*StartBasketResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for StartBasket")
	}

	var r0 *StartBasketResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *StartBasketRequest) (*StartBasketResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *StartBasketRequest) *StartBasketResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*StartBasketResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *StartBasketRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedBasketServiceServer provides a mock function with no fields
func (_m *MockBasketServiceServer) mustEmbedUnimplementedBasketServiceServer() {
	_m.Called()
}

// NewMockBasketServiceServer creates a new instance of MockBasketServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBasketServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBasketServiceServer {
	mock := &MockBasketServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
