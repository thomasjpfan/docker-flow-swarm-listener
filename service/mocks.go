package service

import (
	"github.com/stretchr/testify/mock"
)

type notificationSenderMock struct {
	mock.Mock
}

func (m *notificationSenderMock) Create(params string) error {
	args := m.Called(params)
	return args.Error(0)
}

func (m *notificationSenderMock) Remove(params string) error {
	args := m.Called(params)
	return args.Error(0)
}

func (m *notificationSenderMock) GetCreateAddr() string {
	args := m.Called()
	return args.String(0)
}

func (m *notificationSenderMock) GetRemoveAddr() string {
	args := m.Called()
	return args.String(0)
}
