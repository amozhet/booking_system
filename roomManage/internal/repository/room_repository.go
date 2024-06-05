package repository

import (
	"database/sql"
	"fmt"
	"roomManage/internal/domain/model"
)

type RoomRepository interface {
	GetAllRooms(filter string, sort string, page int, pageSize int) ([]model.Room, error)
	GetRoomByID(id int) (*model.Room, error)
	CreateRoom(room *model.Room) error
	UpdateRoom(room *model.Room) error
	DeleteRoom(id int) error
}

type RoomRepositoryImpl struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &RoomRepositoryImpl{db: db}
}

func (r *RoomRepositoryImpl) GetAllRooms(filter string, sort string, page int, pageSize int) ([]model.Room, error) {
	query := "SELECT id, name, description, available FROM rooms"

	if filter != "" {
		query += " WHERE name ILIKE '%" + filter + "%'"
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Available); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepositoryImpl) GetRoomByID(id int) (*model.Room, error) {
	var room model.Room
	err := r.db.QueryRow("SELECT id, name, description, available FROM rooms WHERE id = $1", id).
		Scan(&room.ID, &room.Name, &room.Description, &room.Available)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepositoryImpl) CreateRoom(room *model.Room) error {
	_, err := r.db.Exec("INSERT INTO rooms (name, description, available) VALUES ($1, $2, $3)",
		room.Name, room.Description, room.Available)
	return err
}

func (r *RoomRepositoryImpl) UpdateRoom(room *model.Room) error {
	_, err := r.db.Exec("UPDATE rooms SET name = $1, description = $2, available = $3 WHERE id = $4",
		room.Name, room.Description, room.Available, room.ID)
	return err
}

func (r *RoomRepositoryImpl) DeleteRoom(id int) error {
	_, err := r.db.Exec("DELETE FROM rooms WHERE id = $1", id)
	return err
}
