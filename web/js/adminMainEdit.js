document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const editButton = document.getElementById("edit-button"); // Кнопка "Изменить"
    let selectedRow = null; // Текущая выбранная строка
    let selectedType = ""; // Тип записи: "client" или "booking" или "general-game"

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
                } else if (parentSection.id === "general-games") {
                    selectedType = "general-game";
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
            const clientName = cells[2]?.textContent.trim(); // Имя клиента
            const tariffName = cells[3]?.textContent.trim(); // имя тарифа
            const generalGamesCell = cells[4]?.textContent.trim(); // Общие игры
            const bookingDate = cells[6]?.textContent.trim(); // Дата бронирования
            const bookingTime = cells[7]?.textContent.trim(); // Время бронирования

            // Разбиваем строку игр на массив (убираем запятые и пробелы)
            const generalGames = generalGamesCell
                ? generalGamesCell.split(",").map(game => game.trim()).filter(game => game)
                : [];

            // Добавляем массив игр в URL
            const generalGamesParam = generalGames.map(game => `games[]=${encodeURIComponent(game)}`).join("&");

            editUrl = `/admin/booking/edit?id=${encodeURIComponent(recordId)}&client=${encodeURIComponent(clientName)}&tariff=${encodeURIComponent(tariffName)}&${generalGamesParam}&date=${encodeURIComponent(bookingDate)}&time=${encodeURIComponent(bookingTime)}`;
        } else if (selectedType === "general-game") {
            const name = cells[1]?.textContent.trim(); // Название игры
            const genre = cells[2]?.textContent.trim(); // Жанр
            const price = cells[3]?.textContent.trim(); // Цена

            editUrl = `/admin/general-game/edit?id=${encodeURIComponent(recordId)}&name=${encodeURIComponent(name)}&genre=${encodeURIComponent(genre)}&price=${encodeURIComponent(price)}`;
        }

        // Переход на страницу редактирования
        if (editUrl) {
            window.location.href = editUrl;
        }
    });
});