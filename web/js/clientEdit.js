document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const editButton = document.getElementById("edit-button"); // Кнопка "Изменить"
    let selectedRow = null; // Текущая выбранная строка

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
                editButton.classList.add("disabled");
                return;
            }

            // Убираем выделение с других строк
            rows.forEach(r => r.classList.remove("selected"));

            // Добавляем выделение для текущей строки
            row.classList.add("selected");
            selectedRow = row;
            editButton.classList.remove("disabled"); // Активируем кнопку "Изменить"
        });
    });

    // Обработка кнопки "Изменить"
    editButton.addEventListener("click", function () {
        if (!selectedRow) return;

        // Получаем данные из выбранной строки
        const cells = selectedRow.querySelectorAll("td");
        const clientId = cells[0]?.textContent.trim(); // ID клиента из скрытого атрибута
        const name = cells[1]?.textContent.trim(); // Имя клиента
        const email = cells[3]?.textContent.trim(); // Email клиента
        const phone = cells[2]?.textContent.trim(); // Телефон клиента

        // Формируем URL для редактирования
        const editUrl = `/admin/client/edit?id=${encodeURIComponent(clientId)}&name=${encodeURIComponent(name)}&email=${encodeURIComponent(email)}&phone=${encodeURIComponent(phone)}`;

        // Переход на страницу редактирования
        window.location.href = editUrl;
    });
});