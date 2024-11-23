package main

import (
	"net/http"

	"VR_project/database"
	"VR_project/internal/interfaces"
)

func main() {
	database.Connect()
	interfaces.HandlerStatic()
	interfaces.HandlerPages()
	http.ListenAndServe(":8080", nil)
}
