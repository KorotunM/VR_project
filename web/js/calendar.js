document.addEventListener("DOMContentLoaded", () => {
    const dateInput = document.getElementById("booking-date");
    const availableTimesSection = document.getElementById("available-times");
    const timeButtonsContainer = document.getElementById("time-buttons");
    const cancelSelectionButton = document.getElementById("cancel-selection");
    const bookingForm = document.getElementById("booking-form");

    let selectedTime = null; // Для хранения выбранного времени

    // Функция выделения времени
    function selectTime(selectedButton, time) {
        document.querySelectorAll(".time-button").forEach(button => {
            button.classList.remove("selected");
            button.classList.add("dimmed");
        });
        selectedButton.classList.add("selected");
        selectedButton.classList.remove("dimmed");

        selectedTime = time;
        cancelSelectionButton.style.display = "inline-block";
    }

    // Запрос времени при выборе даты
    dateInput.addEventListener("change", () => {
        const selectedDate = dateInput.value;

        if (selectedDate) {
            fetch(`/available-times?date=${selectedDate}`)
                .then(response => response.json())
                .then(availableTimes => {
                    availableTimesSection.style.display = "block";

                    timeButtonsContainer.innerHTML = "";
                    availableTimes.forEach(time => {
                        const button = document.createElement("button");
                        button.textContent = time;
                        button.classList.add("time-button");
                        button.type = "button";
                        button.addEventListener("click", () => selectTime(button, time));
                        timeButtonsContainer.appendChild(button);
                    });
                });
        } else {
            availableTimesSection.style.display = "none";
        }
    });

    // Обработка отправки формы
    bookingForm.addEventListener("submit", function (e) {
        e.preventDefault();

        const formData = {
            name: document.getElementById("name").value,
            email: document.getElementById("email").value,
            phone: document.getElementById("phone").value,
            tariff: document.getElementById("tariff-select").value,
            booking_date: document.getElementById("booking-date").value,
            booking_time: selectedTime || "12:00"
        };

        fetch("/booking", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) throw new Error("Ошибка сервера");
            return response.text();
        })
        .then(data => {
            alert("Бронирование успешно отправлено!");
            clearFormAndButtons(bookingForm); // Очистка формы и кнопок
        })
        .catch(error => {
            console.error("Ошибка:", error);
            alert("Произошла ошибка при отправке бронирования.");
        });
    });
});
