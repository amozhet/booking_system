package main

import (
	"clientManage/internal/app"
	"clientManage/internal/repository"
	"clientManage/internal/service"
	grpcTransport "clientManage/internal/transport/grpc"
	"clientManage/internal/transport/http/handler"
	"clientManage/internal/transport/http/middleware"
	"clientManage/internal/transport/messaging"
	"clientManage/pkg/database"
	customLogger "clientManage/pkg/logger"
	pb "clientManage/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	logger := customLogger.NewLogger()

	cfg := app.LoadConfig()

	db, err := database.ConnectDB(cfg.DBUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	clientRepo := repository.NewClientRepository(db)

	clientMessaging, err := messaging.NewClientMessaging(cfg.RabbitMQUrl)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := clientMessaging.Close(); err != nil {
			logger.Fatalf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	clientService := service.NewClientService(clientRepo, clientMessaging)

	grpcServer := grpcTransport.NewClientGRPCServer(clientService)
	grpcSrv := grpc.NewServer()
	pb.RegisterClientManagementServiceServer(grpcSrv, grpcServer)

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

	clientHandler := handler.NewClientHandler(clientService)

	r := mux.NewRouter()

	r.HandleFunc("/register", clientHandler.RegisterClient).Methods("POST")
	r.HandleFunc("/login", clientHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)

	api.Handle("/clients", middleware.RoleMiddleware("client", "admin")(http.HandlerFunc(clientHandler.ListClients))).Methods("GET")
	api.Handle("/clients/{client_id}", middleware.RoleMiddleware("admin")(http.HandlerFunc(clientHandler.GetClient))).Methods("GET")
	api.Handle("/clients", middleware.RoleMiddleware("admin")(http.HandlerFunc(clientHandler.CreateClient))).Methods("POST")
	api.Handle("/clients/{client_id}", middleware.RoleMiddleware("admin")(http.HandlerFunc(clientHandler.UpdateClient))).Methods("PUT")
	api.Handle("/clients/{client_id}", middleware.RoleMiddleware("admin")(http.HandlerFunc(clientHandler.DeleteClient))).Methods("DELETE")

	logger.Println("HTTP server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
