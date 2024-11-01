// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/marvelalexius/jones/model"
	mock "github.com/stretchr/testify/mock"

	stripe "github.com/stripe/stripe-go/v76"
)

// ISubscriptionService is an autogenerated mock type for the ISubscriptionService type
type ISubscriptionService struct {
	mock.Mock
}

// CustomerPortal provides a mock function with given fields: ctx, userID
func (_m *ISubscriptionService) CustomerPortal(ctx context.Context, userID string) (string, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for CustomerPortal")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleInvoicePaid provides a mock function with given fields: ctx, invoice
func (_m *ISubscriptionService) HandleInvoicePaid(ctx context.Context, invoice *stripe.Invoice) error {
	ret := _m.Called(ctx, invoice)

	if len(ret) == 0 {
		panic("no return value specified for HandleInvoicePaid")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *stripe.Invoice) error); ok {
		r0 = rf(ctx, invoice)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleInvoicePaymentFailed provides a mock function with given fields: ctx, customerEmail
func (_m *ISubscriptionService) HandleInvoicePaymentFailed(ctx context.Context, customerEmail string) error {
	ret := _m.Called(ctx, customerEmail)

	if len(ret) == 0 {
		panic("no return value specified for HandleInvoicePaymentFailed")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, customerEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleSubscriptionDeleted provides a mock function with given fields: ctx, stripeSubscription
func (_m *ISubscriptionService) HandleSubscriptionDeleted(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	ret := _m.Called(ctx, stripeSubscription)

	if len(ret) == 0 {
		panic("no return value specified for HandleSubscriptionDeleted")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *stripe.Subscription) error); ok {
		r0 = rf(ctx, stripeSubscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleSubscriptionUpdated provides a mock function with given fields: ctx, stripeSubscription
func (_m *ISubscriptionService) HandleSubscriptionUpdated(ctx context.Context, stripeSubscription *stripe.Subscription) error {
	ret := _m.Called(ctx, stripeSubscription)

	if len(ret) == 0 {
		panic("no return value specified for HandleSubscriptionUpdated")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *stripe.Subscription) error); ok {
		r0 = rf(ctx, stripeSubscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: ctx, userID, req
func (_m *ISubscriptionService) Subscribe(ctx context.Context, userID string, req model.SubscriptionRequest) (string, error) {
	ret := _m.Called(ctx, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for Subscribe")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.SubscriptionRequest) (string, error)); ok {
		return rf(ctx, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, model.SubscriptionRequest) string); ok {
		r0 = rf(ctx, userID, req)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, model.SubscriptionRequest) error); ok {
		r1 = rf(ctx, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewISubscriptionService creates a new instance of ISubscriptionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISubscriptionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISubscriptionService {
	mock := &ISubscriptionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
