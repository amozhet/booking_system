package repository

import (
	"booking/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

type BookingRepositoryMock struct {
	mock.Mock
}

func (m *BookingRepositoryMock) CreateBooking(booking *model.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetBookingByID(id int64) (*model.Booking, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Booking), args.Error(1)
}

func (m *BookingRepositoryMock) UpdateBooking(booking *model.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingRepositoryMock) DeleteBooking(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BookingRepositoryMock) ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error) {
	args := m.Called(offset, limit, filters, sortBy, sortOrder)
	return args.Get(0).([]*model.Booking), args.Error(1)
}
