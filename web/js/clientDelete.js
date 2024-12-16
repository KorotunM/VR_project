document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const deleteButton = document.getElementById("delete-button");
    let selectedRow = null; // Текущая выбранная строка
    let selectedId = ""; // Id выбранного пользователя

    rows.forEach(row => {
        row.addEventListener("click", function () {
            if (row.classList.contains("add-row")) {
                return; // Прерываем обработку, если клик был по строке с кнопкой
            }

            // Если строка уже выбрана, снимаем выделение
            if (selectedRow === row) {
                row.classList.remove("selected");
                selectedRow = null;
                selectedId = "";
                deleteButton.classList.add("disabled");
                return;
            }

            // Убираем выделение со всех строк
            rows.forEach(r => r.classList.remove("selected"));

            // Выделяем текущую строку
            row.classList.add("selected");
            selectedRow = row;
            selectedId = row.querySelector(".hidden-id").textContent.trim();
            deleteButton.classList.remove("disabled");
        });
    });

    // Обработка нажатия кнопки удаления
    deleteButton.addEventListener("click", function () {
        if (!selectedRow || !selectedId) return;

        // Подтверждение удаления
        const confirmDelete = confirm("Вы уверены, что хотите удалить этого клиента?");
        if (!confirmDelete) return;

        // Перенаправление с добавлением Id в URL
        window.location.href = `/admin/client/delete?id=${encodeURIComponent(selectedId)}`;
    });
});