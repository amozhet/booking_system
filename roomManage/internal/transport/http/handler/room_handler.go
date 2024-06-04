package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"roomManage/internal/domain/model"
	"roomManage/internal/service"
	"strconv"
)

type RoomsHandler struct {
	service *service.RoomsService
}

func NewRoomHandler(service *service.RoomsService) *RoomsHandler {
	return &RoomsHandler{service: service}
}

func getUserRole(r *http.Request) string {
	return r.Header.Get("Role")
}

func (h *RoomsHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var rooms model.Room
	err := json.NewDecoder(r.Body).Decode(&rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateRoom(&rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomsHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["room_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	room, err := h.service.GetRoomsByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (h *RoomsHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["room_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	var rooms model.Room
	err = json.NewDecoder(r.Body).Decode(&rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rooms.ID = id

	err = h.service.UpdateRoom(&rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomsHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["room_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteRoom(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RoomsHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
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

	rooms, err := h.service.ListRooms(offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}
