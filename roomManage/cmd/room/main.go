package main

import (
	"roomManage/internal/app"
	"roomManage/internal/config"
	"roomManage/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	application := app.NewApp(cfg, log)

	go func() {
		if err := application.RunHTTPServer(); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	if err := application.RunGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
