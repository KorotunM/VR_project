document.addEventListener("DOMContentLoaded", function () {
    const dateInput = document.getElementById("date");
    const timeSelect = document.getElementById("time");

    // Обновление времени при выборе новой даты
    dateInput.addEventListener("change", function () {
        const selectedDate = dateInput.value;
        if (!selectedDate) return;

        // Очистка селектора времени перед обновлением
        timeSelect.innerHTML = '<option value="">Загрузка...</option>';

        // AJAX-запрос для получения доступных времён
        fetch(`/available-times?date=${encodeURIComponent(selectedDate)}`)
            .then(response => response.json())
            .then(data => {
                timeSelect.innerHTML = "";

                if (data.times && data.times.length > 0) {
                    data.times.forEach(time => {
                        const option = document.createElement("option");
                        option.value = time;
                        option.textContent = time;

                        // Сохраняем выбранное время, если оно уже задано
                        if (time === "{{.BookingTime}}") {
                            option.selected = true;
                        }

                        timeSelect.appendChild(option);
                    });
                } else {
                    const option = document.createElement("option");
                    option.value = "";
                    option.textContent = "Нет доступного времени";
                    timeSelect.appendChild(option);
                }
            })
            .catch(error => {
                console.error("Ошибка при загрузке времени:", error);
                timeSelect.innerHTML = '<option value="">Ошибка загрузки</option>';
            });
    });
});