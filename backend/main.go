package main

import (
	"log"
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/handler"
)

func main() {
	db := database.Connect()
	router := handler.SetupRoutes(db)

	log.Println("Backend running at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}