package interfaces

import (
	"VR_project/internal/services"
	"net/http"
)

func HandlerStatic() {
	// обработка css и js
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../web/"))))
	// обработка изображений
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
}

func HandlerPages() {
	// Отображение страниц
	http.HandleFunc("/", TariffHandler)
	http.HandleFunc("/available-times", AvailableTimesHandler)
	http.HandleFunc("/booking", HandleBooking)
	http.HandleFunc("/admin", AdminPage)
	http.HandleFunc("/admin/tariff", TariffPage)
	http.HandleFunc("/admin/tariff/delete/element", services.DeleteElementTariff)
	http.HandleFunc("/admin/tariff/add/game", AddGamePage)
	http.HandleFunc("/admin/tariff/add/device", AddDevicePage)
	http.HandleFunc("/admin/tariff/edit/game", EditGamePage)
	http.HandleFunc("/admin/tariff/edit/device", EditDevicePage)
	http.HandleFunc("/admin/tariff/add", AddTariffPage)
	http.HandleFunc("/admin/tariff/edit", EditTariffPage)
	http.HandleFunc("/admin/tariff/delete", services.DeleteTariff)
	http.HandleFunc("/admin/client/delete", services.DeleteClient)
	http.HandleFunc("/admin/client/add", AddClientPage)
	http.HandleFunc("/admin/client/edit", EditClientPage)
	http.HandleFunc("/admin/booking/delete", services.DeleteBooking)
	http.HandleFunc("/admin/booking/edit", EditBookingPage)
}
