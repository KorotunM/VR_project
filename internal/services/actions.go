package services

import (
	"VR_project/database"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteElementTariff(w http.ResponseWriter, r *http.Request) {
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
			return "Такая игра уже есть", nil
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
			return "Такое устройство уже есть", nil
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
		return "", fmt.Errorf("error editing device in tariff: %v", err)
	}
	http.Redirect(w, r, "/admin/tariff?id="+tariffId, http.StatusSeeOther)
	return "", nil
}
