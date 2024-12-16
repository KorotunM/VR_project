package interfaces

import (
	"html/template"
	"log"
	"net/http"

	"VR_project/database"
)

// TariffHandler обрабатывает запрос для отображения тарифов
func TariffHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем тарифы из MongoDB
	tariffs, err := database.GetTariffs()
	if err != nil {
		http.Error(w, "Ошибка получения данных о тарифах", http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("../web/templates/Client/index.html")
	if err != nil {
		log.Println("Ошибка загрузки шаблона:", err)
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Передаём данные тарифов в шаблон
	data := map[string]interface{}{
		"Tariffs": tariffs,
	}

	// Рендеринг шаблона с данными
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка обработки шаблона", http.StatusInternalServerError)
	}
}
