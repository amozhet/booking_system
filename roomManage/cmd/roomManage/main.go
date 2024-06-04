package main

import (
	"log"
	"os"
	"roomManage/internal/app"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := app.LoadConfig()
	application := app.NewApplication(cfg, logger)
	application.Run()
}
