package repository

import (
	"errors"
	"roomManage/internal/domain/model"
	"sort"
)

type RoomRepository struct {
	rooms map[string]*model.Room
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*model.Room),
	}
}

func (r *RoomRepository) GetAll() ([]*model.Room, error) {
	var rooms []*model.Room
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepository) GetByID(id string) (*model.Room, error) {
	if room, exists := r.rooms[id]; exists {
		return room, nil
	}
	return nil, errors.New("room not found")
}

func (r *RoomRepository) Save(room *model.Room) error {
	r.rooms[room.ID] = room
	return nil
}

func (r *RoomRepository) Delete(id string) error {
	delete(r.rooms, id)
	return nil
}

func (r *RoomRepository) Filter(predicate func(*model.Room) bool) ([]*model.Room, error) {
	var filteredRooms []*model.Room
	for _, room := range r.rooms {
		if predicate(room) {
			filteredRooms = append(filteredRooms, room)
		}
	}
	return filteredRooms, nil
}

func (r *RoomRepository) Sort(compare func(a, b *model.Room) bool) ([]*model.Room, error) {
	rooms, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	sort.Slice(rooms, func(i, j int) bool {
		return compare(rooms[i], rooms[j])
	})
	return rooms, nil
}

func (r *RoomRepository) Paginate(rooms []*model.Room, page, pageSize int) ([]*model.Room, error) {
	start := (page - 1) * pageSize
	if start >= len(rooms) {
		return nil, nil
	}
	end := start + pageSize
	if end > len(rooms) {
		end = len(rooms)
	}
	return rooms[start:end], nil
}
