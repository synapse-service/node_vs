// Code generated by mockery. DO NOT EDIT.

package service

import (
	context "context"

	google_golang_orggrpc "google.golang.org/grpc"

	grpc "github.com/synapse-service/gateway/transport/grpc"

	mock "github.com/stretchr/testify/mock"
)

// GatewayClient is an autogenerated mock type for the GatewayClient type
type GatewayClient struct {
	mock.Mock
}

// Start provides a mock function with given fields: _a0
func (_m *GatewayClient) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields: _a0
func (_m *GatewayClient) Stop(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Sync provides a mock function with given fields: ctx, in, opts
func (_m *GatewayClient) Sync(ctx context.Context, in *grpc.SyncRequest, opts ...google_golang_orggrpc.CallOption) (grpc.GatewayAPI_SyncClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 grpc.GatewayAPI_SyncClient
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.SyncRequest, ...google_golang_orggrpc.CallOption) grpc.GatewayAPI_SyncClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(grpc.GatewayAPI_SyncClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc.SyncRequest, ...google_golang_orggrpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGatewayClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewGatewayClient creates a new instance of GatewayClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGatewayClient(t mockConstructorTestingTNewGatewayClient) *GatewayClient {
	mock := &GatewayClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}