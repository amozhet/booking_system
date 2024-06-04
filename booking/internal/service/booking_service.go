package service

import (
	"booking/internal/domain/model"
	"booking/internal/repository"
	"booking/internal/transport/mesagging"
	"log"
)

type BookingService struct {
	repo      *repository.BookingRepository
	messaging *messaging.BookingMessaging
}

func NewBookingService(repo *repository.BookingRepository, messaging *messaging.BookingMessaging) *BookingService {
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
