package service

import (
	"roomManage/internal/domain/model"
	"roomManage/internal/repository"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (s *RoomService) GetAllRooms() ([]*model.Room, error) {
	return s.repo.GetAll()
}

func (s *RoomService) GetRoomByID(id string) (*model.Room, error) {
	return s.repo.GetByID(id)
}

func (s *RoomService) CreateRoom(room *model.Room) error {
	return s.repo.Save(room)
}

func (s *RoomService) UpdateRoom(id string, room *model.Room) error {
	return s.repo.Save(room)
}

func (s *RoomService) DeleteRoom(id string) error {
	return s.repo.Delete(id)
}

func (s *RoomService) FilterRooms(predicate func(*model.Room) bool) ([]*model.Room, error) {
	return s.repo.Filter(predicate)
}

func (s *RoomService) SortRooms(compare func(a, b *model.Room) bool) ([]*model.Room, error) {
	return s.repo.Sort(compare)
}

func (s *RoomService) PaginateRooms(page, pageSize int) ([]*model.Room, error) {
	rooms, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.repo.Paginate(rooms, page, pageSize)
}
