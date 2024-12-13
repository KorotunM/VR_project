document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const deleteButton = document.getElementById("delete-button");
    let selectedRow = null; // Текущая выбранная строка
    let selectedType = ""; // Тип: "game" или "device"

     // Добавляем _id тарифа в URL
     const urlParams = new URLSearchParams(window.location.search);
     const tariffId = urlParams.get("id"); // Получаем значение параметра "id"
     if (!tariffId) {
         console.error("Не удалось определить _id тарифа из URL");
         return;
     }
 

    rows.forEach(row => {
        row.addEventListener("click", function () {
            // Определяем таблицу, из которой строка (игры или устройства)
            const parentTable = row.closest("section").id;
            selectedType = parentTable === "games" ? "game" : "device";

            // Если строка уже выбрана, снимаем выделение
            if (selectedRow === row) {
                row.classList.remove("selected");
                selectedRow = null;
                deleteButton.classList.add("disabled");
                return;
            }

            // Убираем выделение со всех строк
            rows.forEach(r => r.classList.remove("selected"));

            // Выделяем текущую строку
            row.classList.add("selected");
            selectedRow = row;
            deleteButton.classList.remove("disabled");
        });
    });

    // Удаление строки
    deleteButton.addEventListener("click", function () {
        if (!selectedRow) return;

        // Сбор данных из выбранной строки
        const name = selectedRow.querySelector("td").textContent.trim(); // Имя игры или устройства

        // AJAX-запрос для удаления
        fetch(`/admin/tariff/delete/element?id=${tariffId}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ type: selectedType, name })
        })
            .then(response => response.json())
            .then(result => {
                if (result.success) {
                    selectedRow.remove(); // Удаляем строку из таблицы
                    deleteButton.classList.add("disabled"); // Делаем кнопку неактивной
                    selectedRow = null;
                } else {
                    alert("Ошибка при удалении: " + result.message);
                }
            })
            .catch(error => {
                console.error("Ошибка запроса:", error);
            });
    });
});
