package repository

import (
	"database/sql"
	"log"
	"roomManage/internal/domain/model"
	"strconv"
	"strings"
)

type RoomRepository struct {
	DB *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{DB: db}
}

func (r *RoomRepository) CreateRoom(rooms *model.Room) error {
	query := `INSERT INTO rooms (name, type, description, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(query, rooms.Name, rooms.Type, rooms.Description, rooms.Status).Scan(&rooms.ID)
	if err != nil {
		log.Printf("Error creating rooms: %v", err)
		return err
	}
	return nil
}

func (r *RoomRepository) GetRoomByID(id int64) (*model.Room, error) {
	query := `SELECT id, name, type, description, status FROM rooms WHERE id = $1`
	var rooms model.Room
	err := r.DB.QueryRow(query, id).Scan(&rooms.ID, &rooms.Name, &rooms.Type, &rooms.Description, &rooms.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting rooms by ID: %v", err)
		return nil, err
	}
	return &rooms, nil
}

func (r *RoomRepository) UpdateRoom(rooms *model.Room) error {
	query := `UPDATE rooms SET name = $1, type = $2, description = $3, status = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, rooms.Name, rooms.Type, rooms.Description, rooms.Status, rooms.ID)
	if err != nil {
		log.Printf("Error updating rooms: %v", err)
		return err
	}
	return nil
}

func (r *RoomRepository) DeleteRoom(id int64) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting rooms: %v", err)
		return err
	}
	return nil
}

func (r *RoomRepository) ListRoom(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Room, error) {
	query := `SELECT id, name, type, description, status FROM rooms`
	var whereClauses []string
	var args []interface{}
	i := 1

	for key, value := range filters {
		whereClauses = append(whereClauses, key+" = $"+strconv.Itoa(i))
		args = append(args, value)
		i++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY " + sortBy + " " + sortOrder
	query += " LIMIT $" + strconv.Itoa(i) + " OFFSET $" + strconv.Itoa(i+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error listing rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []*model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Type, &room.Description, &room.Status); err != nil {
			log.Printf("Error scanning room: %v", err)
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error with rows in list rooms: %v", err)
		return nil, err
	}
	return rooms, nil
}
