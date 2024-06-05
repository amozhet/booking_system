package main

import (
	"log"
	"net/http"

	"roomManage/internal/app"
)

func main() {
	a := app.New()
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
