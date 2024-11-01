// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	model "github.com/marvelalexius/jones/model"
	mock "github.com/stretchr/testify/mock"
)

// INotificationRepository is an autogenerated mock type for the INotificationRepository type
type INotificationRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: notif
func (_m *INotificationRepository) Create(notif model.Notification) error {
	ret := _m.Called(notif)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(model.Notification) error); ok {
		r0 = rf(notif)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewINotificationRepository creates a new instance of INotificationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewINotificationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *INotificationRepository {
	mock := &INotificationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
