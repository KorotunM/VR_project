document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const editButton = document.getElementById("edit-button"); // Кнопка "Изменить"
    let selectedRow = null; // Текущая выбранная строка
    let selectedType = ""; // Тип записи: "client" или "booking"

    // Обработка выбора строки таблицы
    rows.forEach(row => {
        row.addEventListener("click", function () {
            // Пропускаем строку с кнопкой "Добавить"
            if (row.classList.contains("add-row")) {
                return;
            }

            // Снимаем выделение, если строка уже выбрана
            if (selectedRow === row) {
                row.classList.remove("selected");
                selectedRow = null;
                selectedType = "";
                editButton.classList.add("disabled");
                return;
            }

            // Убираем выделение с других строк
            rows.forEach(r => r.classList.remove("selected"));

            // Определяем тип записи по родительскому блоку
            const parentSection = row.closest("section");
            if (parentSection) {
                if (parentSection.id === "clients") {
                    selectedType = "client";
                } else if (parentSection.id === "bookings") {
                    selectedType = "booking";
                }
            }

            // Добавляем выделение для текущей строки
            row.classList.add("selected");
            selectedRow = row;
            editButton.classList.remove("disabled"); // Активируем кнопку "Изменить"
        });
    });

    // Обработка кнопки "Изменить"
    editButton.addEventListener("click", function () {
        if (!selectedRow || !selectedType) return;

        // Получаем данные из выбранной строки
        const cells = selectedRow.querySelectorAll("td");
        const recordId = cells[0]?.textContent.trim(); // ID записи из скрытого атрибута

        // Формируем URL для редактирования в зависимости от типа записи
        let editUrl = "";
        if (selectedType === "client") {
            const name = cells[1]?.textContent.trim(); // Имя клиента
            const email = cells[3]?.textContent.trim(); // Email клиента
            const phone = cells[2]?.textContent.trim(); // Телефон клиента

            editUrl = `/admin/client/edit?id=${encodeURIComponent(recordId)}&name=${encodeURIComponent(name)}&email=${encodeURIComponent(email)}&phone=${encodeURIComponent(phone)}`;
        } else if (selectedType === "booking") {
            const bookingDate = cells[3]?.textContent.trim(); // Дата бронирования
            const clientName = cells[1]?.textContent.trim(); // Имя клиента
            const tariffName = cells[2]?.textContent.trim(); // имя тарифа
            const bookingTime = cells[4]?.textContent.trim(); // Время бронирования

            editUrl = `/admin/booking/edit?id=${encodeURIComponent(recordId)}&client=${encodeURIComponent(clientName)}&tariff=${encodeURIComponent(tariffName)}&date=${encodeURIComponent(bookingDate)}&time=${encodeURIComponent(bookingTime)}`;
        }

        // Переход на страницу редактирования
        if (editUrl) {
            window.location.href = editUrl;
        }
    });
});