package interfaces

import (
	"VR_project/database"
	"VR_project/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

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
	bookings, err := database.GetAllBookings()
	if err != nil {
		fmt.Fprintf(w, "Error receiving bookings: %v", err)
		return
	}
	adminPageData.Statistic, err = database.GetBookingStatistics()
	if err != nil {
		fmt.Fprintf(w, "Error getting statistic: %v", err)
		return
	}
	adminPageData.StatisticDays, err = database.GetDailyBookingStatistics()
	if err != nil {
		fmt.Fprintf(w, "Error getting statistic days: %v", err)
		return
	}
	tmp, err := template.ParseFiles("../web/templates/admin/admin.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template: %v", err)
		return
	}
	adminPageData.Clients = clients
	adminPageData.Tariffs = tariffs
	adminPageData.Bookings = bookings

	// заполнение имен клиентов и названий тарифов
	for i := range adminPageData.Bookings {
		for j := range adminPageData.Clients {
			if adminPageData.Bookings[i].ClientID == adminPageData.Clients[j].Id {
				adminPageData.Bookings[i].ClientName = adminPageData.Clients[j].Name
			}
		}
		for j := range adminPageData.Tariffs {
			if adminPageData.Bookings[i].TariffID == adminPageData.Tariffs[j].Id {
				adminPageData.Bookings[i].TariffName = adminPageData.Tariffs[j].Name
			}
		}
	}

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
		answer               database.AdminFormTariff
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
		answer                            database.AdminFormTariff
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
		answer               database.AdminFormTariff
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
		answer                               database.AdminFormTariff
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

func AddClientPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer database.AdminFormClient
		err    error
	)
	if r.Method == http.MethodPost {
		err = services.AddClient(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error adding client: %v", err)
			return
		}
	}
	answer.Action = "Добавить"
	tmp, err := template.ParseFiles("../web/templates/admin/formClient.html")
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

func EditClientPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer             database.AdminFormClient
		name, email, phone string
		err                error
	)

	name = r.URL.Query().Get("name")
	email = r.URL.Query().Get("email")
	phone = r.URL.Query().Get("phone")

	if name == "" || email == "" || phone == "" {
		fmt.Fprintf(w, "Error getting parameters from URL")
		return
	}

	answer.Action = "Редактировать"
	answer.Name = name
	answer.Email = email
	answer.Phone = phone

	if r.Method == http.MethodPost {
		err = services.EditClient(w, r)
		if err != nil {
			fmt.Fprintf(w, "Error editing client: %v", err)
			return
		}
	}
	tmp, err := template.ParseFiles("../web/templates/admin/formClient.html")
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

func EditBookingPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer               database.AdminFormBooking
		client, tariff, date string
		time                 string
		err                  error
		tariffs              []database.TariffTitle
		clients              []database.Client
	)

	client = r.URL.Query().Get("client")
	tariff = r.URL.Query().Get("tariff")
	date = r.URL.Query().Get("date")
	time = r.URL.Query().Get("time")

	answer.Action = "Редактировать"
	answer.ClientName = client
	answer.TariffName = tariff
	answer.BookingDate = date
	answer.BookingTime = time

	tariffs, err = database.GetAllTariffs()
	if err != nil {
		fmt.Fprintf(w, "Error getting all tariffs: %v", err)
		return
	}

	answer.Tariffs = tariffs

	clients, err = database.GetClients()
	if err != nil {
		fmt.Fprintf(w, "Error getting all clients: %v", err)
		return
	}

	answer.Clients = clients

	answer.AvailableTimes = []string{"10:00", "12:00", "14:00", "16:00", "18:00", "20:00"}

	if r.Method == http.MethodPost {
		err = services.EditBooking(w, r)
		if err != nil {
			if err.Error() == "time already exist" {
				// Если ошибка валидации (занятое время), заполняем поле ошибки
				answer.Validation = "Это время уже занято"
			} else {
				fmt.Fprintf(w, "Error editing booking: %v", err)
				return
			}
		}
	}

	tmp, err := template.ParseFiles("../web/templates/admin/formBooking.html")
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

func AddBookingPage(w http.ResponseWriter, r *http.Request) {
	var (
		answer database.AdminFormBooking
		err    error
	)

	tariffs, err := database.GetAllTariffs()
	if err != nil {
		fmt.Fprintf(w, "Error getting all tariffs: %v", err)
		return
	}

	answer.Tariffs = tariffs

	clients, err := database.GetClients()
	if err != nil {
		fmt.Fprintf(w, "Error getting all clients: %v", err)
		return
	}

	answer.Clients = clients

	answer.AvailableTimes = []string{"10:00", "12:00", "14:00", "16:00", "18:00", "20:00"}

	if r.Method == http.MethodPost {
		err = services.AddBooking(w, r)
		if err != nil {
			if err.Error() == "time already exist" {
				// Если ошибка валидации (занятое время), заполняем поле ошибки
				answer.Validation = "Это время уже занято"
			} else {
				fmt.Fprintf(w, "Error adding booking: %v", err)
				return
			}
		}
	}
	answer.Action = "Добавить"
	tmp, err := template.ParseFiles("../web/templates/admin/formBooking.html")
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

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// Здесь будет храниться логин и зашифрованный пароль администратора
		adminUsername     = "admin"
		adminPasswordHash = "$2a$10$JZCb33DSy2qcbTiEUMn8XeOM0jjCFgsDxlhKtE6aCdbMhTxpU/ovG"
	)

	// Проверяем, если метод запроса POST
	if r.Method == "POST" {
		// Получаем данные из формы
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка логина и пароля
		if username == adminUsername {
			err := bcrypt.CompareHashAndPassword([]byte(adminPasswordHash), []byte(password))
			if err == nil {
				// Если логин и пароль совпали, перенаправляем на админку
				http.Redirect(w, r, "/admin", http.StatusFound)
				return
			}

		}
	}

	tmp, err := template.ParseFiles("../web/templates/admin/formAdmin.html")
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
