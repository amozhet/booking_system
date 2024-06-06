package unit

import (
	"reflect"
	"testing"

	"roomManage/internal/domain/model"
	"roomManage/internal/repository"
	"roomManage/internal/service"
)

func setup() *service.RoomService {
	roomRepo := repository.NewMockRoomRepository()
	return service.NewRoomService((*repository.RoomRepository)(roomRepo))
}

func TestCreateRoom(t *testing.T) {
	svc := setup()

	room := &model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}

	err := svc.CreateRoom(room)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	storedRoom, err := svc.GetRoomByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(storedRoom, room) {
		t.Fatalf("Expected %v, got %v", room, storedRoom)
	}
}

func TestGetRoomByID(t *testing.T) {
	svc := setup()

	room := &model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}

	svc.CreateRoom(room)
	storedRoom, err := svc.GetRoomByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(storedRoom, room) {
		t.Fatalf("Expected %v, got %v", room, storedRoom)
	}
}

func TestUpdateRoom(t *testing.T) {
	svc := setup()

	room := &model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}

	svc.CreateRoom(room)

	updatedRoom := &model.Room{
		ID:          "1",
		Name:        "Updated Room 1",
		Description: "Updated Description 1",
		Available:   false,
	}

	err := svc.UpdateRoom("1", updatedRoom)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	storedRoom, err := svc.GetRoomByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(storedRoom, updatedRoom) {
		t.Fatalf("Expected %v, got %v", updatedRoom, storedRoom)
	}
}

func TestDeleteRoom(t *testing.T) {
	svc := setup()

	room := &model.Room{
		ID:          "1",
		Name:        "Room 1",
		Description: "Description 1",
		Available:   true,
	}

	svc.CreateRoom(room)

	err := svc.DeleteRoom("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = svc.GetRoomByID("1")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGetAllRooms(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	storedRooms, err := svc.GetAllRooms()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(storedRooms, rooms) {
		t.Fatalf("Expected %v, got %v", rooms, storedRooms)
	}
}

func TestFilterRooms(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	predicate := func(room *model.Room) bool {
		return room.Available
	}

	filteredRooms, err := svc.FilterRooms(predicate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRooms := []*model.Room{rooms[0]}
	if !reflect.DeepEqual(filteredRooms, expectedRooms) {
		t.Fatalf("Expected %v, got %v", expectedRooms, filteredRooms)
	}
}

func TestSortRooms(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	compare := func(a, b *model.Room) bool {
		return a.Name < b.Name
	}

	sortedRooms, err := svc.SortRooms(compare)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRooms := []*model.Room{rooms[1], rooms[0]}
	if !reflect.DeepEqual(sortedRooms, expectedRooms) {
		t.Fatalf("Expected %v, got %v", expectedRooms, sortedRooms)
	}
}

func TestPaginateRooms(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
		{ID: "3", Name: "Room 3", Description: "Description 3", Available: true},
		{ID: "4", Name: "Room 4", Description: "Description 4", Available: false},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	paginatedRooms, err := svc.PaginateRooms(1, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRooms := []*model.Room{rooms[0], rooms[1]}
	if !reflect.DeepEqual(paginatedRooms, expectedRooms) {
		t.Fatalf("Expected %v, got %v", expectedRooms, paginatedRooms)
	}
}

func TestPaginateRooms_Page2(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
		{ID: "3", Name: "Room 3", Description: "Description 3", Available: true},
		{ID: "4", Name: "Room 4", Description: "Description 4", Available: false},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	paginatedRooms, err := svc.PaginateRooms(2, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRooms := []*model.Room{rooms[2], rooms[3]}
	if !reflect.DeepEqual(paginatedRooms, expectedRooms) {
		t.Fatalf("Expected %v, got %v", expectedRooms, paginatedRooms)
	}
}

func TestSortAndPaginateRooms(t *testing.T) {
	svc := setup()

	rooms := []*model.Room{
		{ID: "3", Name: "Room 3", Description: "Description 3", Available: true},
		{ID: "1", Name: "Room 1", Description: "Description 1", Available: true},
		{ID: "4", Name: "Room 4", Description: "Description 4", Available: false},
		{ID: "2", Name: "Room 2", Description: "Description 2", Available: false},
	}

	for _, room := range rooms {
		svc.CreateRoom(room)
	}

	compare := func(a, b *model.Room) bool {
		return a.Name < b.Name
	}

	sortedRooms, err := svc.SortRooms(compare)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	paginatedRooms, err := svc.PaginateRooms(1, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRooms := []*model.Room{sortedRooms[0], sortedRooms[1]}
	if !reflect.DeepEqual(paginatedRooms, expectedRooms) {
		t.Fatalf("Expected %v, got %v", expectedRooms, paginatedRooms)
	}
}
