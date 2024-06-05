package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"roomManage/internal/config"
	"roomManage/internal/repository"
	"roomManage/internal/service"
	handler "roomManage/internal/transport/http"
	"roomManage/pkg/database"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func New() *App {
	config := config.LoadConfig()
	db, err := database.NewDB(config.DBUrl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	router := mux.NewRouter()
	router.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods("GET")
	router.HandleFunc("/rooms/{room_id}", roomHandler.GetRoomByID).Methods("GET")
	router.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	router.HandleFunc("/rooms/{room_id}", roomHandler.UpdateRoom).Methods("PUT")
	router.HandleFunc("/rooms/{room_id}", roomHandler.DeleteRoom).Methods("DELETE")

	return &App{
		Router: router,
		DB:     db,
	}
}
