package messaging

import (
	"clientManage/internal/data"
	"github.com/stretchr/testify/mock"
)

type UserMessagingMock struct {
	mock.Mock
}

func (m *UserMessagingMock) PublishClientCreated(user *data.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserMessagingMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
