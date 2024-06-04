package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"roomManage/internal/config"
	"roomManage/internal/repository"
	"roomManage/internal/service"
	"roomManage/internal/transport/http/handler"
	"roomManage/internal/transport/mesagging"
	"roomManage/pkg/database"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := app.LoadConfig()

	db, err := database.ConnectDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	roomRepo := repository.NewRoomRepository(db)
	roomMessaging, err := messaging.NewRoomMessaging(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer roomMessaging.Close()

	roomService := service.NewRoomService(roomRepo, roomMessaging)
	roomHandler := handler.NewRoomHandler(roomService)

	r := mux.NewRouter()
	r.HandleFunc("/rooms", roomHandler.ListRooms).Methods("GET")
	r.HandleFunc("/rooms/{room_id}", roomHandler.GetRoom).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}", roomHandler.UpdateRoom).Methods("PUT")
	r.HandleFunc("/rooms/{room_id}", roomHandler.DeleteRoom).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	logger.Printf("Server is running on port %s", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
