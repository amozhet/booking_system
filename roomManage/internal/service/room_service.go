package service

import (
	"roomManage/internal/domain/model"
	"roomManage/internal/repository"
)

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

/*func (s *RoomService) GetAllRooms(string, string, int, int) ([]model.Room, error) {
	return s.repo.GetAllRooms()
}*/

func (s *RoomService) GetRoomByID(id int) (*model.Room, error) {
	return s.repo.GetRoomByID(id)
}

func (s *RoomService) CreateRoom(room *model.Room) error {
	return s.repo.CreateRoom(room)
}

func (s *RoomService) UpdateRoom(room *model.Room) error {
	return s.repo.UpdateRoom(room)
}

func (s *RoomService) DeleteRoom(id int) error {
	return s.repo.DeleteRoom(id)
}
