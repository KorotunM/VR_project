document.addEventListener("DOMContentLoaded", function () {
    const bookingRows = document.querySelectorAll("#bookings table tbody tr");
    const clientRows = document.querySelectorAll("#clients table tbody tr");
    const tariffButtons = document.querySelectorAll("#tariffs .tariff-button");
    const generalGamesRows = document.querySelectorAll("#general-games table tbody tr");

    // Класс для выделения
    const highlightClass = "highlight";

    // Переменная для хранения последнего выбранного элемента бронирования
    let lastSelectedBooking = null;

    // Снятие выделения со всех строк клиентов, тарифов и игр
    function clearHighlight() {
        clientRows.forEach(row => row.classList.remove(highlightClass));
        tariffButtons.forEach(button => button.classList.remove(highlightClass));
        generalGamesRows.forEach(row => row.classList.remove(highlightClass));
    }

    // Обработка клика по строке бронирования
    bookingRows.forEach(row => {
        row.addEventListener("click", function () {
            if (row.classList.contains("add-row")) {
                return; // Пропускаем строку "Добавить"
            }

            // Если это повторный клик по той же строке
            if (lastSelectedBooking === row) {
                clearHighlight(); // Снимаем выделение
                lastSelectedBooking = null; // Сбрасываем текущий выбор
                return;
            }

            // Устанавливаем новую строку как последнюю выбранную
            lastSelectedBooking = row;

            // Снимаем выделение перед выделением новых элементов
            clearHighlight();

            // Получаем ID клиента, имя тарифа и массив игр
            const clientId = row.querySelector(".hidden-id:nth-child(2)")?.textContent.trim();
            const tariffName = row.querySelector(".table-body-cell:nth-child(4)")?.textContent.trim();
            const gamesList = row.querySelector(".table-body-cell:nth-child(5)")?.textContent.trim();
            const gamesArray = gamesList ? gamesList.split(',').map(game => game.trim()) : [];

            // Выделяем соответствующую строку клиента
            if (clientId) {
                clientRows.forEach(clientRow => {
                    const clientRowId = clientRow.querySelector(".hidden-id")?.textContent.trim();
                    if (clientRowId === clientId) {
                        clientRow.classList.add(highlightClass);
                    }
                });
            }

            // Выделяем соответствующий тариф
            if (tariffName) {
                tariffButtons.forEach(button => {
                    if (button.textContent.trim() === tariffName) {
                        button.classList.add(highlightClass);
                    }
                });
            }

            // Выделяем соответствующие игры
            if (gamesArray.length > 0) {
                generalGamesRows.forEach(gameRow => {
                    const gameName = gameRow.querySelector(".table-body-cell:nth-child(2)")?.textContent.trim();
                    if (gamesArray.includes(gameName)) {
                        gameRow.classList.add(highlightClass);
                    }
                });
            }
        });
    });

    // Обработка клика по строке клиента
    clientRows.forEach(row => {
        row.addEventListener("click", function () {
            // Снимаем предыдущее выделение
            clearHighlight();
        });
    });

    // Обработка клика по строке игры
    generalGamesRows.forEach(row => {
        row.addEventListener("click", function () {
            // Снимаем предыдущее выделение
            clearHighlight();
        });
    });
});
