package main

import (
	"clientManage/internal/app"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := app.LoadConfig()
	application := app.NewApplication(cfg, logger)
	application.Run()
}
