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

    // Функция для очистки формы
    function clearForm() {
        if (bookingForm) {
            bookingForm.reset(); // Сбрасывает все поля формы
            console.log("Форма сброшена.");
        }

        selectedTime = null; // Сбрасываем выбранное время

        // Скрываем кнопку "Отменить выбор", если она есть
        if (cancelSelectionButton) {
            cancelSelectionButton.style.display = "none";
        }

        // Сбрасываем состояние кнопок времени
        clearTimeButtons();
    }

    // Слушатель события для кнопки "Забронировать" (после успешной отправки)
    document.addEventListener("formSubmitted", () => {
        clearForm();
    });
});
