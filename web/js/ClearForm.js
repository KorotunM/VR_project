// ClearForm.js: функция для очистки формы и кнопок времени

function clearFormAndButtons(bookingForm) {
    const timeButtonsContainer = document.getElementById("time-buttons");
    const cancelSelectionButton = document.getElementById("cancel-selection");

    // Очистка формы
    if (bookingForm) {
        bookingForm.reset();
        console.log("Форма успешно очищена.");
    }

    // Сброс выбранного времени
    selectedTime = null;

    // Очистка стилей кнопок времени
    if (timeButtonsContainer) {
        document.querySelectorAll(".time-button").forEach(button => {
            button.classList.remove("selected", "dimmed");
        });
    }

    // Скрываем кнопку отмены выбора
    if (cancelSelectionButton) {
        cancelSelectionButton.style.display = "none";
    }
}
