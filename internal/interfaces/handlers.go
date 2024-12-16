package interfaces

import (
	"VR_project/database"
	"VR_project/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	tmp, err := template.ParseFiles("../web/templates/index.html")
// 	if err != nil {
// 		fmt.Fprintf(w, "Error loading template: %v", err)
// 		return
// 	}
// 	err = tmp.Execute(w, nil)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error rendering template: %v", err)
// 		return
// 	}
// }

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
		answer               database.AdminFormTariffData
		validation, tariffId string
		err                  error
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		fmt.Fprintf(w, "Error getting id tariff from URL")
		return
	}
	answer.IdTariff = tariffId
	answer.Action = "Добавить"
	if r.Method == http.MethodPost {
		validation, err = services.AddGame(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding game: %v", err)
			return
		}
		answer.Validation = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formGame.html")
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

func EditGamePage(w http.ResponseWriter, r *http.Request) {
	var (
		answer                            database.AdminFormTariffData
		validation, tariffId, name, genre string
		err                               error
	)
	tariffId = r.URL.Query().Get("id")
	name = r.URL.Query().Get("name")
	genre = r.URL.Query().Get("genre")

	if tariffId == "" || name == "" || genre == "" {
		fmt.Fprintf(w, "Error getting parameters from URL")
		return
	}
	answer.IdTariff = tariffId
	answer.Action = "Редактировать"
	answer.Name = name
	answer.Genre = genre

	if r.Method == http.MethodPost {
		validation, err = services.EditGame(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error editing game: %v", err)
			return
		}
		answer.Validation = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formGame.html")
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
		answer               database.AdminFormTariffData
		validation, tariffId string
		err                  error
	)
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		fmt.Fprintf(w, "Error getting id tariff from URL")
		return
	}
	answer.IdTariff = tariffId
	answer.Action = "Добавить"
	if r.Method == http.MethodPost {
		validation, err = services.AddDevice(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding device: %v", err)
			return
		}
		answer.Validation = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formDevice.html")
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

func EditDevicePage(w http.ResponseWriter, r *http.Request) {
	var (
		answer                               database.AdminFormTariffData
		validation, tariffId, name, platform string
		err                                  error
	)
	tariffId = r.URL.Query().Get("id")
	name = r.URL.Query().Get("name")
	platform = r.URL.Query().Get("platform")

	if tariffId == "" || name == "" || platform == "" {
		fmt.Fprintf(w, "Error getting parameters from URL")
		return
	}

	answer.IdTariff = tariffId
	answer.Action = "Редактировать"
	answer.Name = name
	answer.Platform = platform

	if r.Method == http.MethodPost {
		validation, err = services.EditDevice(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error editing device: %v", err)
			return
		}
		answer.Validation = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formDevice.html")
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

func AddTariffPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer     database.AdminFormTariff
		validation string
		err        error
	)
	if r.Method == http.MethodPost {
		validation, err = services.AddTariff(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding tariff: %v", err)
			return
		}
		answer.Validation = validation
	}
	answer.Action = "Добавить"
	tmp, err := template.ParseFiles("../web/templates/admin/formTariff.html")
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

func EditTariffPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer                     database.AdminFormTariff
		validation, name, tariffId string
		err                        error
		price                      int
	)

	name = r.URL.Query().Get("name")
	tariffId = r.URL.Query().Get("id")
	price, err = strconv.Atoi(r.URL.Query().Get("price"))

	if name == "" || tariffId == "" {
		fmt.Fprintf(w, "Error getting string parameters from URL")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Error getting tariff price from URL")
		return
	}

	answer.IdTariff = tariffId
	answer.Action = "Редактировать"
	answer.Name = name
	answer.Price = price

	if r.Method == http.MethodPost {
		validation, err = services.EditTariff(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error editing tariff: %v", err)
			return
		}
		answer.Validation = validation
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formTariff.html")
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
