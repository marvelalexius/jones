// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	model "github.com/marvelalexius/jones/model"
	mock "github.com/stretchr/testify/mock"
)

// NotificationRepository is an autogenerated mock type for the NotificationRepository type
type NotificationRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: notif
func (_m *NotificationRepository) Create(notif model.Notification) error {
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

// NewNotificationRepository creates a new instance of NotificationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotificationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotificationRepository {
	mock := &NotificationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
