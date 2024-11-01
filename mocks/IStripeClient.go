// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	stripe "github.com/stripe/stripe-go/v76"
)

// IStripeClient is an autogenerated mock type for the IStripeClient type
type IStripeClient struct {
	mock.Mock
}

// CreateBillingPortalSession provides a mock function with given fields: ctx, customerID
func (_m *IStripeClient) CreateBillingPortalSession(ctx context.Context, customerID string) (*stripe.BillingPortalSession, error) {
	ret := _m.Called(ctx, customerID)

	if len(ret) == 0 {
		panic("no return value specified for CreateBillingPortalSession")
	}

	var r0 *stripe.BillingPortalSession
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*stripe.BillingPortalSession, error)); ok {
		return rf(ctx, customerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *stripe.BillingPortalSession); ok {
		r0 = rf(ctx, customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stripe.BillingPortalSession)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCheckoutSession provides a mock function with given fields: ctx, customerID, planID
func (_m *IStripeClient) CreateCheckoutSession(ctx context.Context, customerID string, planID string) (*stripe.CheckoutSession, error) {
	ret := _m.Called(ctx, customerID, planID)

	if len(ret) == 0 {
		panic("no return value specified for CreateCheckoutSession")
	}

	var r0 *stripe.CheckoutSession
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*stripe.CheckoutSession, error)); ok {
		return rf(ctx, customerID, planID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *stripe.CheckoutSession); ok {
		r0 = rf(ctx, customerID, planID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stripe.CheckoutSession)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, customerID, planID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCustomer provides a mock function with given fields: ctx, email, name
func (_m *IStripeClient) CreateCustomer(ctx context.Context, email string, name string) (*stripe.Customer, error) {
	ret := _m.Called(ctx, email, name)

	if len(ret) == 0 {
		panic("no return value specified for CreateCustomer")
	}

	var r0 *stripe.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*stripe.Customer, error)); ok {
		return rf(ctx, email, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *stripe.Customer); ok {
		r0 = rf(ctx, email, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stripe.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIStripeClient creates a new instance of IStripeClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIStripeClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IStripeClient {
	mock := &IStripeClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}