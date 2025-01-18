document.addEventListener("DOMContentLoaded", () => {
    const bookingForm = document.getElementById("booking-form");
    const cancelSelectionButton = document.getElementById("cancel-selection");
    const timeButtonsContainer = document.getElementById("time-buttons");

    let selectedTime = null;

    // Функция для сброса кнопок времени
    function clearTimeButtons() {
        document.querySelectorAll(".time-button").forEach(button => {
            button.classList.remove("selected", "dimmed");
        });
        console.log("Кнопки времени сброшены.");
    }

    // Функция для сброса только поля даты и области времени
    function clearDateAndTime() {
        // Очищаем поле с датой
        const dateField = document.getElementById("booking-date");
        if (dateField) {
            dateField.value = ""; // Сбрасываем значение поля даты
            console.log("Поле даты очищено.");
        }

        selectedTime = null; // Сбрасываем выбранное время

        // Скрываем кнопку "Отменить выбор"
        if (cancelSelectionButton) {
            cancelSelectionButton.style.display = "none";
        }

        // Скрываем область доступного времени
        const availableTimes = document.getElementById("available-times");
        if (availableTimes) {
            availableTimes.style.display = "none";
        }

        // Сбрасываем кнопки времени
        clearTimeButtons();
    }

    // Добавляем слушатель на кнопку "Отменить выбор"
    if (cancelSelectionButton) {
        cancelSelectionButton.addEventListener("click", () => {
            clearDateAndTime();
        });
    }

    // Слушатель события для кнопки "Забронировать" (после успешной отправки)
    document.addEventListener("formSubmitted", () => {
        clearDateAndTime();
    });
});
