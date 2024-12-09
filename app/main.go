package main

import (
	"net/http"

	"VR_project/database"
	"VR_project/internal/interfaces"
)

func main() {
	// Инициализация MongoDB
	database.InitMongoDB("mongodb://localhost:27017")
	defer database.CloseMongoDB()

	interfaces.HandlerStatic()
	interfaces.HandlerPages()
	http.ListenAndServe(":8080", nil)
}
