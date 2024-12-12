package interfaces

import (
	"VR_project/database"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../web/templates/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
		return
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	var adminPageData database.AdminPageData
	clients, err := database.GetClients()
	if err != nil {
		fmt.Fprintf(w, "Error receiving clients: %v", err)
		return
	}
	tariffs, err := database.GetAllTariffs()
	if err != nil {
		fmt.Fprintf(w, "Error receiving tariffs: %v", err)
		return
	}
	tmp, err := template.ParseFiles("../web/templates/admin.html")
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

func TariffPage(w http.ResponseWriter, r *http.Request) {
	var (
		tariff database.Tariff
		err    error
	)
	tariff, err = database.GetTariff(r)
	if err != nil {
		fmt.Fprintf(w, "Error getting the tariff: %v", err)
		return
	}
	tmpl, err := template.ParseFiles("../web/templates/tariff.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}

	err = tmpl.Execute(w, tariff)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
		return
	}
}

func DeleteElementTariff(w http.ResponseWriter, r *http.Request) {
	var err error = database.DeleteElementTariffDB(r)
	if err != nil {
		fmt.Fprintf(w, "Error deleting element tariff: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}
