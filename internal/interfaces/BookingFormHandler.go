package interfaces

import (
	"VR_project/database"
	"encoding/json"
	"log"
	"net/http"
)

func HandleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req database.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Шаг 1: Записываем клиента в Clients
	clientID, err := database.InsertClient(req.Name, req.Email, req.Phone)
	if err != nil {
		log.Printf("Error inserting client: %v", err)
		http.Error(w, "Error inserting client", http.StatusInternalServerError)
		return
	}

	// Шаг 2: Получаем tariff_id по имени тарифа
	tariffName := getTariffName(req.Tariff) // Функция обрезает текст до тире
	tariffID, err := database.GetTariffIDByName(tariffName)
	if err != nil {
		log.Printf("Tariff not found: %v", err)
		http.Error(w, "Tariff not found", http.StatusBadRequest)
		return
	}

	// Шаг 3: Записываем бронирование в Bookings
	err = database.InsertBooking(clientID, tariffID, req.BookingDate, req.BookingTime)
	if err != nil {
		log.Printf("Error inserting booking: %v", err)
		http.Error(w, "Error inserting booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Booking successful"))
}

func getTariffName(fullName string) string {
	// Обрезаем имя тарифа до первого тире
	for i, char := range fullName {
		if char == '-' {
			return fullName[:i-1]
		}
	}
	return fullName
}
