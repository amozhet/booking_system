package service

import (
	"booking/internal/domain/model"
	"booking/internal/repository"
	"booking/internal/transport/messaging"
	"log"
)

type BookingService struct {
	repo      repository.BookingRepository
	messaging messaging.BookingMessaging
}

func NewBookingService(repo repository.BookingRepository, messaging messaging.BookingMessaging) *BookingService {
	return &BookingService{repo: repo, messaging: messaging}
}

func (s *BookingService) CreateBooking(booking *model.Booking) error {
	err := s.repo.CreateBooking(booking)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return err
	}

	err = s.messaging.PublishBookingCreated(booking)
	if err != nil {
		log.Printf("Error publishing booking created message: %v", err)
		return err
	}

	return nil
}

func (s *BookingService) GetBookingByID(id int64) (*model.Booking, error) {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		log.Printf("Error getting booking by ID: %v", err)
		return nil, err
	}
	return booking, nil
}

func (s *BookingService) UpdateBooking(booking *model.Booking) error {
	err := s.repo.UpdateBooking(booking)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		return err
	}
	return nil
}

func (s *BookingService) DeleteBooking(id int64) error {
	err := s.repo.DeleteBooking(id)
	if err != nil {
		log.Printf("Error deleting booking: %v", err)
		return err
	}
	return nil
}

func (s *BookingService) ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error) {
	bookings, err := s.repo.ListBookings(offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		log.Printf("Error listing bookings: %v", err)
		return nil, err
	}
	return bookings, nil
}
