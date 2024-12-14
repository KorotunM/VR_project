package interfaces

import (
	"VR_project/database"
	"VR_project/internal/services"
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
	tmp, err := template.ParseFiles("../web/templates/admin/admin.html")
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
	tmpl, err := template.ParseFiles("../web/templates/admin/tariff.html")
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

func AddGamePage(w http.ResponseWriter, r *http.Request) {
	var (
		answer               database.Validation
		validation, tariffId string
		err                  error
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		fmt.Fprintf(w, "Error getting id tariff from URL")
		return
	}
	answer.IdTariff = tariffId
	if r.Method == http.MethodPost {
		validation, err = services.AddGame(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding game: %v", err)
			return
		}
		answer.Error = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/addGame.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}
	err = tmp.Execute(w, answer)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
		return
	}
}

func AddDevicePage(w http.ResponseWriter, r *http.Request) {
	var (
		answer               database.Validation
		validation, tariffId string
		err                  error
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		fmt.Fprintf(w, "Error getting id tariff from URL")
		return
	}
	answer.IdTariff = tariffId
	if r.Method == http.MethodPost {
		validation, err = services.AddDevice(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding device: %v", err)
			return
		}
		answer.Error = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/addDevice.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}
	err = tmp.Execute(w, answer)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
		return
	}
}
