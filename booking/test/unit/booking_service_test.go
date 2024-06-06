package service_test

import (
	"booking/internal/domain/model"
	"booking/internal/repository"
	"booking/internal/service"
	"booking/internal/transport/messaging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() (*repository.BookingRepositoryMock, *messaging.BookingMessagingMock, *service.BookingService) {
	repoMock := new(repository.BookingRepositoryMock)
	messagingMock := new(messaging.BookingMessagingMock)
	svc := service.NewBookingService(repoMock, messagingMock)
	return repoMock, messagingMock, svc
}

func TestCreateBooking(t *testing.T) {
	repoMock, messagingMock, svc := setup()

	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}

	repoMock.On("CreateBooking", booking).Return(nil)
	messagingMock.On("PublishBookingCreated", booking).Return(nil)

	err := svc.CreateBooking(booking)
	assert.Nil(t, err)
	assert.NotZero(t, booking.ID)

	repoMock.AssertCalled(t, "CreateBooking", booking)
	messagingMock.AssertCalled(t, "PublishBookingCreated", booking)
}

func TestGetBookingByID(t *testing.T) {
	repoMock, _, svc := setup()

	booking := &model.Booking{
		ID:        1,
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}

	repoMock.On("GetBookingByID", int64(1)).Return(booking, nil)
	result, err := svc.GetBookingByID(1)
	assert.Nil(t, err)
	assert.Equal(t, booking, result)
	repoMock.AssertExpectations(t)
}

func TestUpdateBooking(t *testing.T) {
	repoMock, _, svc := setup()

	booking := &model.Booking{
		ID:        1,
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}

	repoMock.On("UpdateBooking", booking).Return(nil)
	err := svc.UpdateBooking(booking)
	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestDeleteBooking(t *testing.T) {
	repoMock, _, svc := setup()

	repoMock.On("DeleteBooking", int64(1)).Return(nil)

	err := svc.DeleteBooking(1)
	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestListBookings(t *testing.T) {
	repoMock, _, svc := setup()

	bookings := []*model.Booking{
		{
			ID:        1,
			ClientID:  1,
			RoomID:    1,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			Status:    "confirmed",
		},
	}

	filters := make(map[string]interface{})
	repoMock.On("ListBookings", 0, 10, filters, "id", "asc").Return(bookings, nil)

	result, err := svc.ListBookings(0, 10, filters, "id", "asc")
	assert.Nil(t, err)
	assert.Equal(t, bookings, result)
	repoMock.AssertExpectations(t)
}

func TestCreateBookingWithError(t *testing.T) {
	repoMock, messagingMock, svc := setup()

	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}

	repoMock.On("CreateBooking", booking).Return(assert.AnError)
	messagingMock.On("PublishBookingCreated", booking).Return(nil)

	err := svc.CreateBooking(booking)
	assert.NotNil(t, err)
	repoMock.AssertExpectations(t)
	messagingMock.AssertExpectations(t)
}

func TestGetBookingByIDNotFound(t *testing.T) {
	repoMock, _, svc := setup()

	repoMock.On("GetBookingByID", int64(1)).Return(nil, assert.AnError)

	result, err := svc.GetBookingByID(1)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	repoMock.AssertExpectations(t)
}

func TestUpdateBookingWithError(t *testing.T) {
	repoMock, _, svc := setup()

	booking := &model.Booking{
		ID:        1,
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}

	repoMock.On("UpdateBooking", booking).Return(assert.AnError)

	err := svc.UpdateBooking(booking)
	assert.NotNil(t, err)
	repoMock.AssertExpectations(t)
}

func TestDeleteBookingWithError(t *testing.T) {
	repoMock, _, svc := setup()

	repoMock.On("DeleteBooking", int64(1)).Return(assert.AnError)

	err := svc.DeleteBooking(1)
	assert.NotNil(t, err)
	repoMock.AssertExpectations(t)
}

func TestListBookingsWithError(t *testing.T) {
	repoMock, _, svc := setup()

	filters := make(map[string]interface{})
	repoMock.On("ListBookings", 0, 10, filters, "id", "asc").Return(nil, assert.AnError)

	result, err := svc.ListBookings(0, 10, filters, "id", "asc")
	assert.NotNil(t, err)
	assert.Nil(t, result)
	repoMock.AssertExpectations(t)
}
