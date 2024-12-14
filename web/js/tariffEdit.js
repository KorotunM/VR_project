document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const editButton = document.getElementById("edit-button"); // Кнопка "Изменить"
    let selectedRow = null; // Текущая выбранная строка
    let selectedType = ""; // Тип: "game" или "device"

    // Получаем ID тарифа из URL
    const urlParams = new URLSearchParams(window.location.search);
    const tariffId = urlParams.get("id"); // id тарифа

    if (!tariffId) {
        console.error("ID тарифа не найден в URL");
        return;
    }

    // Обработка выбора строки таблицы
    rows.forEach(row => {
        row.addEventListener("click", function () {
            // Пропускаем строку с кнопкой "Добавить"
            if (row.classList.contains("add-row")) {
                return;
            }

            // Определяем тип: игра или устройство
            const parentTable = row.closest("section").id;
            selectedType = parentTable === "games" ? "game" : "device";

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
        const name = cells[0]?.textContent.trim(); // Название (первая колонка)

        let editUrl = "";

        if (selectedType === "game") {
            const genre = cells[1]?.textContent.trim(); // Жанр
            editUrl = `/admin/tariff/edit/game?id=${encodeURIComponent(tariffId)}&name=${encodeURIComponent(name)}&genre=${encodeURIComponent(genre)}`;
        } else if (selectedType === "device") {
            const platform = cells[1]?.textContent.trim(); // Платформа
            editUrl = `/admin/tariff/edit/device?id=${encodeURIComponent(tariffId)}&name=${encodeURIComponent(name)}&platform=${encodeURIComponent(platform)}`;
        }

        // Переход на страницу редактирования
        window.location.href = editUrl;
    });
});