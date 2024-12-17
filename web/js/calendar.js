document.addEventListener("DOMContentLoaded", () => {
    const dateInput = document.getElementById("booking-date");
    const availableTimesSection = document.getElementById("available-times");
    const timeButtonsContainer = document.getElementById("time-buttons");
    const cancelSelectionButton = document.getElementById("cancel-selection");
    const bookingForm = document.getElementById("booking-form");

    let selectedTime = null; // Для хранения выбранного времени

    // Функция для отображения кнопок времени
    function displayTimeButtons(times) {
        timeButtonsContainer.innerHTML = ""; // Очищаем контейнер

        times.forEach(time => {
            const button = document.createElement("button");
            button.textContent = time;
            button.classList.add("time-button");
            button.type = "button";

            // Обработчик нажатия кнопки времени
            button.addEventListener("click", () => {
                selectTime(button, time);
            });

            timeButtonsContainer.appendChild(button);
        });

        // Скрываем кнопку отмены выбора
        cancelSelectionButton.style.display = "none";
    }

    // Функция выделения времени
    function selectTime(selectedButton, time) {
        // Сбрасываем стиль всех кнопок
        document.querySelectorAll(".time-button").forEach(button => {
            button.classList.remove("selected");
            button.classList.add("dimmed");
        });

        // Выделяем выбранную кнопку
        selectedButton.classList.add("selected");
        selectedButton.classList.remove("dimmed");

        // Сохраняем выбранное время
        selectedTime = time;
        console.log("Выбрано время:", selectedTime);

        // Показываем кнопку "Отменить выбор"
        cancelSelectionButton.style.display = "inline-block";
    }

    // Функция отмены выбора
    cancelSelectionButton.addEventListener("click", () => {
        // Сбрасываем стиль всех кнопок
        document.querySelectorAll(".time-button").forEach(button => {
            button.classList.remove("selected", "dimmed");
        });

        // Сбрасываем выбранное время
        selectedTime = null;
        console.log("Выбор времени отменен");

        // Скрываем кнопку "Отменить выбор"
        cancelSelectionButton.style.display = "none";
    });

    // Запрос времени при выборе даты
    dateInput.addEventListener("change", () => {
        const selectedDate = dateInput.value;

        if (selectedDate) {
            fetch(`/available-times?date=${selectedDate}`)
                .then(response => response.json())
                .then(availableTimes => {
                    availableTimesSection.style.display = "block";

                    if (availableTimes.length > 0) {
                        displayTimeButtons(availableTimes);
                    } else {
                        timeButtonsContainer.innerHTML = "<p>Нет доступного времени</p>";
                    }
                })
                .catch(error => {
                    console.error("Ошибка при получении данных:", error);
                    availableTimesSection.style.display = "block";
                    timeButtonsContainer.innerHTML = "<p>Ошибка загрузки данных</p>";
                });
        } else {
            availableTimesSection.style.display = "none";
        }
    });

    // Обработка отправки формы
    bookingForm.addEventListener("submit", function (e) {
        e.preventDefault(); // Предотвращаем стандартное поведение формы

        const formData = {
            name: document.getElementById("name").value,
            email: document.getElementById("email").value,
            phone: document.getElementById("phone").value,
            tariff: document.getElementById("tariff-select").value,
            booking_date: document.getElementById("booking-date").value,
            booking_time: selectedTime || "12:00" // Используем выбранное время или дефолтное
        };

        console.log("Отправляем данные:", formData);

        fetch("/booking", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) throw new Error("Ошибка сервера");
            return response.text();
        })
        .then(data => {
            console.log("Успешно:", data);
            alert("Бронирование успешно отправлено!");

            document.dispatchEvent(new Event("formSubmitted"));
        })
        .catch(error => {
            console.error("Ошибка:", error);
            alert("Произошла ошибка при отправке бронирования.");
        });
    });
});
