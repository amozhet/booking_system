package handler

import (
	"booking/internal/domain/model"
	"booking/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func getUserRole(r *http.Request) string {
	return r.Header.Get("Role")
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var booking model.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateBooking(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["book_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["book_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var booking model.Booking
	err = json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	booking.ID = id

	err = h.service.UpdateBooking(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["book_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBooking(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	filters := make(map[string]interface{})
	for key, values := range query {
		if key != "offset" && key != "limit" && key != "sort_by" && key != "sort_order" {
			filters[key] = values[0]
		}
	}
	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := query.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	bookings, err := h.service.ListBookings(offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}
