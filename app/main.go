package main

import (
	"net/http"

	"VR_project/internal/interfaces"
)

func main() {
	interfaces.HandlerStatic()
	interfaces.HandlerPages()
	http.ListenAndServe(":8080", nil)
}
