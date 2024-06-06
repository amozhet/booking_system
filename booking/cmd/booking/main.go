package main

import (
	"booking/internal/app"
	"booking/internal/repository"
	"booking/internal/service"
	grpcTransport "booking/internal/transport/grpc"
	"booking/internal/transport/http/handler"
	"booking/internal/transport/messaging"
	"booking/pkg/database"
	customLogger "booking/pkg/logger"
	pb "booking/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	logger := customLogger.NewLogger()

	// Load configuration
	cfg := app.LoadConfig()

	// Connect to the database
	db, err := database.ConnectDB(cfg.DBUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repository
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize messaging
	bookingMessaging, err := messaging.NewBookingMessaging(cfg.RabbitMQUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := bookingMessaging.Close(); err != nil {
			logger.Fatalf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	// Initialize service
	bookingService := service.NewBookingService(bookingRepo, bookingMessaging)

	// Initialize gRPC server
	grpcServer := grpcTransport.NewBookingGRPCServer(bookingService)
	grpcSrv := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcSrv, grpcServer)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatalf("Failed to listen on port 50051: %v", err)
		}
		logger.Println("gRPC server is running on port 50051")
		if err := grpcSrv.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Initialize handler
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/bookings", bookingHandler.ListBookings).Methods("GET")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.GetBooking).Methods("GET")
	r.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.UpdateBooking).Methods("PUT")
	r.HandleFunc("/bookings/{book_id}", bookingHandler.DeleteBooking).Methods("DELETE")

	// Set up and start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	logger.Printf("HTTP server is running on port %s", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
