package messaging

import (
	"booking/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

type BookingMessagingMock struct {
	mock.Mock
}

func (m *BookingMessagingMock) PublishBookingCreated(booking *model.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingMessagingMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
