package interfaces

import (
	"VR_project/database"
	"fmt"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../web/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	tmp.Execute(w, nil)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	var adminPageData database.AdminPageData
	clients, err := database.GetClients()
	if err != nil {
		fmt.Fprintf(w, "Error receiving clients: %v", err)
		return
	}
	tariffs, err := database.GetTariffs()
	if err != nil {
		fmt.Fprintf(w, "Error receiving tariffs: %v", err)
		return
	}
	tmp, err := template.ParseFiles("../web/admin.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}
	adminPageData.Clients = clients
	adminPageData.Tariffs = tariffs
	err = tmp.Execute(w, adminPageData)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
		return
	}
}
