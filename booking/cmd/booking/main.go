package main

import (
	"booking/internal/app"
	"booking/internal/repository"
	"booking/internal/service"
	"booking/internal/transport/http/handler"
	"booking/internal/transport/mesagging"
	"booking/pkg/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := app.LoadConfig()

	db, err := database.ConnectDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bookingRepo := repository.NewBookingRepository(db)
	bookingMessaging, err := messaging.NewBookingMessaging(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer bookingMessaging.Close()

	bookingService := service.NewBookingService(bookingRepo, bookingMessaging)
	bookingHandler := handler.NewBookingHandler(bookingService)

	r := mux.NewRouter()
	r.HandleFunc("/bookings", bookingHandler.ListBookings).Methods("GET")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.GetBooking).Methods("GET")
	r.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.UpdateBooking).Methods("PUT")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.DeleteBooking).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	logger.Printf("Server is running on port %s", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
