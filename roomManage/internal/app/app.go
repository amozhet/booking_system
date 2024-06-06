package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"roomManage/internal/config"
	"roomManage/internal/repository"
	"roomManage/internal/service"
	grpcHandler "roomManage/internal/transport/grpc"
	httpHandler "roomManage/internal/transport/http"
	"roomManage/pkg/logger"
	"roomManage/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type App struct {
	Config *config.Config
	Logger *logger.Logger
}
type Config struct {
	Port        string
	Env         string
	DBUrl       string
	RabbitMQUrl string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "4002"),
		Env:         getEnv("ENV", "development"),
		DBUrl:       getEnv("DB_URL", "postgres://user:pass@localhost/bookingdb"),
		RabbitMQUrl: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: %s environment variable not set. Using default value: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func NewApp(cfg *config.Config, log *logger.Logger) *App {
	return &App{
		Config: cfg,
		Logger: log,
	}
}

func (a *App) RunHTTPServer() error {
	r := mux.NewRouter()
	roomRepo := repository.NewRoomRepository()
	roomService := service.NewRoomService(roomRepo)
	roomHandler := httpHandler.NewRoomHandler(roomService)

	r.HandleFunc("/rooms", roomHandler.GetRooms).Methods("GET")
	r.HandleFunc("/rooms/{room_id}", roomHandler.GetRoomByID).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}", roomHandler.UpdateRoom).Methods("PUT")
	r.HandleFunc("/rooms/{room_id}", roomHandler.DeleteRoom).Methods("DELETE")

	serverAddr := fmt.Sprintf(":%d", a.Config.Server.Port)
	return http.ListenAndServe(serverAddr, r)
}

func (a *App) RunGRPCServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Config.GRPC.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	roomService := service.NewRoomService(repository.NewRoomRepository())
	roomGRPCServer := grpcHandler.NewRoomGRPCServer(roomService)

	proto.RegisterRoomServiceServer(grpcServer, roomGRPCServer)

	return grpcServer.Serve(lis)
}
