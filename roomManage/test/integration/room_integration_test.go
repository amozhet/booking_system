package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"roomManage/internal/domain/model"
	"roomManage/internal/repository"
	"roomManage/internal/service"
	httpHandler "roomManage/internal/transport/http"
)

func setupRouter() *mux.Router {
	roomRepo := repository.NewRoomRepository()
	roomService := service.NewRoomService(roomRepo)
	roomHandler := httpHandler.NewRoomHandler(roomService)

	r := mux.NewRouter()
	r.HandleFunc("/rooms", roomHandler.GetRooms).Methods("GET")
	r.HandleFunc("/rooms/{room_id}", roomHandler.GetRoomByID).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}", roomHandler.UpdateRoom).Methods("PUT")
	r.HandleFunc("/rooms/{room_id}", roomHandler.DeleteRoom).Methods("DELETE")

	return r
}

func TestCreateRoom(t *testing.T) {
	router := setupRouter()

	room := model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}
	body, _ := json.Marshal(room)

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, response.Code)
	}
}

func TestGetRoom(t *testing.T) {
	router := setupRouter()

	room := model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}
	roomRepo := repository.NewRoomRepository()
	roomRepo.Save(&room)

	req, _ := http.NewRequest("GET", "/rooms/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}

	var returnedRoom model.Room
	json.NewDecoder(response.Body).Decode(&returnedRoom)
	if returnedRoom.ID != room.ID {
		t.Errorf("Expected room ID %s, but got %s", room.ID, returnedRoom.ID)
	}
}

func TestUpdateRoom(t *testing.T) {
	router := setupRouter()

	room := model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}
	roomRepo := repository.NewRoomRepository()
	roomRepo.Save(&room)

	updatedRoom := model.Room{
		ID:          "1",
		Name:        "Updated Room 1",
		Description: "Updated Description 1",
		Available:   false,
	}
	body, _ := json.Marshal(updatedRoom)

	req, _ := http.NewRequest("PUT", "/rooms/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("GET", "/rooms/1", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	var returnedRoom model.Room
	json.NewDecoder(response.Body).Decode(&returnedRoom)
	if returnedRoom.Name != updatedRoom.Name {
		t.Errorf("Expected room name %s, but got %s", updatedRoom.Name, returnedRoom.Name)
	}
}

func TestDeleteRoom(t *testing.T) {
	router := setupRouter()

	room := model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}
	roomRepo := repository.NewRoomRepository()
	roomRepo.Save(&room)

	req, _ := http.NewRequest("DELETE", "/rooms/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("GET", "/rooms/1", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, response.Code)
	}
}

func TestGetRooms(t *testing.T) {
	router := setupRouter()

	rooms := []model.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
	}
	roomRepo := repository.NewRoomRepository()
	for _, room := range rooms {
		roomRepo.Save(&room)
	}

	req, _ := http.NewRequest("GET", "/rooms", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}

	var returnedRooms []model.Room
	json.NewDecoder(response.Body).Decode(&returnedRooms)
	if len(returnedRooms) != len(rooms) {
		t.Errorf("Expected %d rooms, but got %d", len(rooms), len(returnedRooms))
	}
}
