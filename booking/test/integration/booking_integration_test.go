package integration_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"booking/internal/app"
	"booking/internal/domain/model"
	"booking/internal/repository"
	"booking/internal/service"
	"booking/internal/transport/http/handler"
	"booking/internal/transport/messaging"
	"booking/pkg/database"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	server           *httptest.Server
	bookingRepo      repository.BookingRepository
	bookingMessaging messaging.BookingMessaging
)

func setup() {
	cfg := app.LoadConfig()

	db, err := database.ConnectDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bookingRepo = repository.NewBookingRepository(db)
	bookingMessaging, err = messaging.NewBookingMessaging(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	bookingService := service.NewBookingService(bookingRepo, bookingMessaging)
	bookingHandler := handler.NewBookingHandler(bookingService)

	r := mux.NewRouter()
	r.HandleFunc("/bookings", bookingHandler.ListBookings).Methods("GET")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.GetBooking).Methods("GET")
	r.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.UpdateBooking).Methods("PUT")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.DeleteBooking).Methods("DELETE")

	server = httptest.NewServer(r)
}

func teardown() {
	server.Close()
	bookingMessaging.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCreateBookingIntegration(t *testing.T) {
	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	body, _ := json.Marshal(booking)

	resp, err := http.Post(server.URL+"/bookings", "application/json", bytes.NewBuffer(body))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdBooking model.Booking
	json.NewDecoder(resp.Body).Decode(&createdBooking)
	assert.NotZero(t, createdBooking.ID)
}

func TestGetBookingByIDIntegration(t *testing.T) {
	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	bookingRepo.CreateBooking(booking)

	resp, err := http.Get(server.URL + "/bookings/" + strconv.FormatInt(booking.ID, 10))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedBooking model.Booking
	json.NewDecoder(resp.Body).Decode(&fetchedBooking)
	assert.Equal(t, booking.ID, fetchedBooking.ID)
}

func TestUpdateBookingIntegration(t *testing.T) {
	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	bookingRepo.CreateBooking(booking)

	booking.EndDate = time.Now().Add(48 * time.Hour)
	body, _ := json.Marshal(booking)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/bookings/"+strconv.FormatInt(booking.ID, 10), bytes.NewBuffer(body))
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedBooking model.Booking
	json.NewDecoder(resp.Body).Decode(&updatedBooking)
	assert.Equal(t, booking.EndDate, updatedBooking.EndDate)
}

func TestDeleteBookingIntegration(t *testing.T) {
	booking := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	bookingRepo.CreateBooking(booking)

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/bookings/"+strconv.FormatInt(booking.ID, 10), nil)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	fetchedBooking, _ := bookingRepo.GetBookingByID(booking.ID)
	assert.Nil(t, fetchedBooking)
}

func TestListBookingsIntegration(t *testing.T) {
	booking1 := &model.Booking{
		ClientID:  1,
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	bookingRepo.CreateBooking(booking1)

	booking2 := &model.Booking{
		ClientID:  2,
		RoomID:    2,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		Status:    "confirmed",
	}
	bookingRepo.CreateBooking(booking2)

	resp, err := http.Get(server.URL + "/bookings")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var bookings []model.Booking
	json.NewDecoder(resp.Body).Decode(&bookings)
	assert.Len(t, bookings, 2)
}
