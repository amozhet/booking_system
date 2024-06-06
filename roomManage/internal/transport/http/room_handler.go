package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"roomManage/internal/domain/model"
	"roomManage/internal/service"
)

type RoomHandler struct {
	service *service.RoomService
}

func NewRoomHandler(service *service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// Пример фильтрации: только доступные комнаты
	predicate := func(room *model.Room) bool {
		return room.Available
	}
	rooms, err := h.service.FilterRooms(predicate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Пример сортировки: по имени комнаты
	compare := func(a, b *model.Room) bool {
		return a.Name < b.Name
	}
	rooms, err = h.service.SortRooms(compare)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Пример пагинации: 1 страница, 10 комнат на странице
	page, pageSize := 1, 10
	rooms, err = h.service.PaginateRooms(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["room_id"]
	room, err := h.service.GetRoomByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room model.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateRoom(&room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *RoomHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["room_id"]
	var room model.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateRoom(id, &room); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["room_id"]
	if err := h.service.DeleteRoom(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
