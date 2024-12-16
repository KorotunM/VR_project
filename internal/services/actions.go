package services

import (
	"VR_project/database"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteElementTariff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "error method")
		return
	}
	var err error = database.DeleteElementTariffDB(r)
	if err != nil {
		fmt.Fprintf(w, "Error deleting element tariff: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}

func AddGame(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		tariffId string
		err      error
	)
	tariffId, err = database.AddGameDB(r)
	if err != nil {
		if err.Error() == "game with this name already exists" {
			return "Такая игра уже существует", nil
		}
		return "", fmt.Errorf("error adding game in tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}

func EditGame(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		tariffId string
		err      error
	)
	tariffId, err = database.EditGameDB(r)
	if err != nil {
		if err.Error() == "game with this name already exists" {
			return "Такая игра уже существует", nil
		}
		return "", fmt.Errorf("error editing game in tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}

func AddDevice(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		tariffId string
		err      error
	)
	tariffId, err = database.AddDeviceDB(r)
	if err != nil {
		if err.Error() == "device with this name already exists" {
			return "Такое устройство уже существует", nil
		}
		return "", fmt.Errorf("error adding device in tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}

func EditDevice(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		tariffId string
		err      error
	)
	tariffId, err = database.EditDeviceDB(r)
	if err != nil {
		if err.Error() == "device with this name already exists" {
			return "Такое устройство уже существует", nil
		}
		return "", fmt.Errorf("error editing device in tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}

func AddTariff(w http.ResponseWriter, r *http.Request) (string, error) {
	var err error = database.AddTariffDB(r)
	if err != nil {
		if err.Error() == "tariff with this name already exists" {
			return "Тариф с таким названием уже существует", nil
		} else if err.Error() == "wrong price" {
			return "Неправильная цена", nil
		}
		return "", fmt.Errorf("error adding tariff: %v", err)
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	return "", nil
}

func EditTariff(w http.ResponseWriter, r *http.Request) (string, error) {
	tariffId, err := database.EditTariffDB(r)
	if err != nil {
		if err.Error() == "tariff with this name already exists" {
			return "Тариф с таким именем уже существует", nil
		}
		return "", fmt.Errorf("error editing tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}

func DeleteTariff(w http.ResponseWriter, r *http.Request) {
	var (
		tariffId string
		err      error
	)

	// Получаем ID тарифа из URL
	tariffId = r.URL.Query().Get("id")
	if tariffId == "" {
		fmt.Fprintf(w, "Error: missing tariff ID")
		return
	}

	// Вызываем функцию удаления тарифа из базы данных
	err = database.DeleteTariffDB(tariffId)
	if err != nil {
		fmt.Fprintf(w, "Error deleting tariff: %v", err)
		return
	}

	// Редирект на главную страницу админки
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Invalid request method")
		return
	}
	err = database.DeleteClientDB(r)
	if err != nil {
		fmt.Fprintf(w, "Error deleting client: %v", err)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
