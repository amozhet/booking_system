package messaging

import (
	"clientManage/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

type ClientMessagingMock struct {
	mock.Mock
}

func (m *ClientMessagingMock) PublishClientCreated(client *model.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientMessagingMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
