package interfaces

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"VR_project/database"
)

// AvailableTimesHandler обрабатывает запросы на доступное время
func AvailableTimesHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Дата не указана", http.StatusBadRequest)
		return
	}

	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	// Получаем забронированные слоты
	bookedSlots, err := database.GetBookedTimes(bookingDate)
	if err != nil {
		http.Error(w, "Ошибка получения данных из базы", http.StatusInternalServerError)
		return
	}

	log.Println("Забронированные слоты из базы данных:", bookedSlots)

	allSlots := []string{"10:00", "12:00", "14:00", "16:00", "18:00", "20:00"}
	availableSlots := []string{}

	for _, slot := range allSlots {
		//log.Println(bookedSlots)
		if !contains(bookedSlots, slot) {
			availableSlots = append(availableSlots, slot)
		}
	}

	log.Println("Доступные временные слоты:", availableSlots)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableSlots)
}

// Проверка наличия временного интервала в списке

func contains(slots []string, slot string) bool {
	for _, s := range slots {
		if strings.TrimSpace(s) == slot {
			return true
		}
	}
	return false
}
