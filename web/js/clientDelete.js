document.addEventListener("DOMContentLoaded", function () {
    const rows = document.querySelectorAll("table tbody tr");
    const deleteButton = document.getElementById("delete-button");
    let selectedRow = null; // Текущая выбранная строка
    let selectedId = ""; // Id выбранной записи
    let selectedType = ""; // Тип записи: "client" или "booking"

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
                selectedType = "";
                deleteButton.classList.add("disabled");
                return;
            }

            // Убираем выделение со всех строк
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

            // Выделяем текущую строку
            row.classList.add("selected");
            selectedRow = row;
            selectedId = row.querySelector(".hidden-id").textContent.trim();
            deleteButton.classList.remove("disabled");
        });
    });

    // Обработка нажатия кнопки удаления
    deleteButton.addEventListener("click", function () {
        if (!selectedRow || !selectedId || !selectedType) return;

        // Подтверждение удаления
        const confirmDelete = confirm("Вы уверены, что хотите удалить эту запись?");
        if (!confirmDelete) return;

        // Формируем URL на основе типа записи
        const baseUrl = selectedType === "client" ? "/admin/client/delete" : "/admin/booking/delete";
        window.location.href = `${baseUrl}?id=${encodeURIComponent(selectedId)}`;
    });
});
