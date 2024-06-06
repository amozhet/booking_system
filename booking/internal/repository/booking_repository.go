package repository

import (
	"booking/internal/domain/model"
	"database/sql"
	"strconv"
	"strings"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int64) (*model.Booking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int64) error
	ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error)
}

type BookingRepositoryImpl struct {
	DB *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepositoryImpl {
	return &BookingRepositoryImpl{DB: db}
}

func (r *BookingRepositoryImpl) CreateBooking(booking *model.Booking) error {
	query := `INSERT INTO bookings (client_id, room_id, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.DB.QueryRow(query, booking.ClientID, booking.RoomID, booking.StartDate, booking.EndDate, booking.Status).Scan(&booking.ID)
}

func (r *BookingRepositoryImpl) GetBookingByID(id int64) (*model.Booking, error) {
	query := `SELECT id, client_id, room_id, start_date, end_date, status FROM bookings WHERE id = $1`
	var booking model.Booking
	err := r.DB.QueryRow(query, id).Scan(&booking.ID, &booking.ClientID, &booking.RoomID, &booking.StartDate, &booking.EndDate, &booking.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepositoryImpl) UpdateBooking(booking *model.Booking) error {
	query := `UPDATE bookings SET client_id = $1, room_id = $2, start_date = $3, end_date = $4, status = $5 WHERE id = $6`
	_, err := r.DB.Exec(query, booking.ClientID, booking.RoomID, booking.StartDate, booking.EndDate, booking.Status, booking.ID)
	return err
}

func (r *BookingRepositoryImpl) DeleteBooking(id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *BookingRepositoryImpl) ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error) {
	query := `SELECT id, client_id, room_id, start_date, end_date, status FROM bookings`
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
		return nil, err
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		var booking model.Booking
		if err := rows.Scan(&booking.ID, &booking.ClientID, &booking.RoomID, &booking.StartDate, &booking.EndDate, &booking.Status); err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}
