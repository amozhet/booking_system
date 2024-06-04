package service

import (
	"booking/internal/domain/model"
	"booking/internal/repository"
	"log"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(booking *model.Booking) error {
	err := s.repo.CreateBooking(booking)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return err
	}
	return nil
}

func (s *BookingService) GetBookingByID(id int64) (*model.Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *BookingService) UpdateBooking(booking *model.Booking) error {
	return s.repo.UpdateBooking(booking)
}

func (s *BookingService) DeleteBooking(id int64) error {
	return s.repo.DeleteBooking(id)
}

func (s *BookingService) ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error) {
	return s.repo.ListBookings(offset, limit, filters, sortBy, sortOrder)
}
