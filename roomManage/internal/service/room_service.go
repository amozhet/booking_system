package service

import (
	"log"
	"roomManage/internal/domain/model"
	"roomManage/internal/repository"
	messaging "roomManage/internal/transport/mesagging"
)

type RoomsService struct {
	repo      *repository.RoomRepository
	messaging *messaging.RoomMessaging
}

func NewRoomService(repo *repository.RoomRepository, messaging *messaging.RoomMessaging) *RoomsService {
	return &RoomsService{repo: repo, messaging: messaging}
}

func (s *RoomsService) CreateRoom(rooms *model.Room) error {
	err := s.repo.CreateRoom(rooms)
	if err != nil {
		log.Printf("Error creating rooms: %v", err)
		return err
	}

	err = s.messaging.PublishRoomCreated(rooms)
	if err != nil {
		log.Printf("Error publishing rooms created message: %v", err)
		return err
	}

	return nil
}
