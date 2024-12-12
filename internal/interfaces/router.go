package interfaces

import "net/http"

func HandlerStatic() {
	// обработка css и js
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../web/"))))
	// обработка изображений
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
}

func HandlerPages() {
	// Отображение страниц
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/admin", AdminPage)
	http.HandleFunc("/admin/tariff", TariffPage)
	http.HandleFunc("/admin/tariff/delete/element", DeleteElementTariff)
}
