package repository

import (
	"booking/internal/domain/model"
	"database/sql"
	"log"
	"strconv"
	"strings"
)

type BookingRepository struct {
	DB *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

func (r *BookingRepository) CreateBooking(booking *model.Booking) error {
	query := `INSERT INTO bookings (client_id, room_id, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.DB.QueryRow(query, booking.ClientID, booking.RoomID, booking.StartDate, booking.EndDate, booking.Status).Scan(&booking.ID)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return err
	}
	return nil
}

func (r *BookingRepository) GetBookingByID(id int64) (*model.Booking, error) {
	query := `SELECT id, client_id, room_id, start_date, end_date, status FROM bookings WHERE id = $1`
	var booking model.Booking
	err := r.DB.QueryRow(query, id).Scan(&booking.ID, &booking.ClientID, &booking.RoomID, &booking.StartDate, &booking.EndDate, &booking.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting booking by ID: %v", err)
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) UpdateBooking(booking *model.Booking) error {
	query := `UPDATE bookings SET client_id = $1, room_id = $2, start_date = $3, end_date = $4, status = $5 WHERE id = $6`
	_, err := r.DB.Exec(query, booking.ClientID, booking.RoomID, booking.StartDate, booking.EndDate, booking.Status, booking.ID)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		return err
	}
	return nil
}

func (r *BookingRepository) DeleteBooking(id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting booking: %v", err)
		return err
	}
	return nil
}

func (r *BookingRepository) ListBookings(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Booking, error) {
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
		log.Printf("Error listing bookings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []*model.Booking
	for rows.Next() {
		var booking model.Booking
		if err := rows.Scan(&booking.ID, &booking.ClientID, &booking.RoomID, &booking.StartDate, &booking.EndDate, &booking.Status); err != nil {
			log.Printf("Error scanning booking: %v", err)
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error with rows in list bookings: %v", err)
		return nil, err
	}
	return bookings, nil
}
