package model

import "time"

type Booking struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	RoomID    int64     `json:"room_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}
