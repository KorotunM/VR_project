document.addEventListener("DOMContentLoaded", function() {
    const tariffSelect = document.getElementById('tariff');
    const selectedTariffId = tariffSelect.value;

    // Если тариф уже выбран при загрузке страницы, обрабатываем его сразу
    if (selectedTariffId) {
        const dropdownList = document.getElementById('dropdown-general-games');
        fetch(`/admin/booking/edit/general-games?tariffId=${encodeURIComponent(selectedTariffId)}`)
            .then(response => response.json())
            .then(selectedGames => {
                const selectedGameIds = new Set(selectedGames.map(game => game.Id));

                // Прячем или показываем игры в зависимости от того, есть ли их Id в ответе от сервера
                Array.from(dropdownList.querySelectorAll('label')).forEach(label => {
                    const input = label.querySelector('input');
                    const gameId = input.value;

                    if (selectedGameIds.has(gameId)) {
                        // Если игра в ответе от сервера, показываем её
                        label.style.display = '';
                        input.disabled = false;
                    } else {
                        // Если игры нет в ответе от сервера, скрываем её
                        label.style.display = 'none';
                        input.disabled = true;
                        input.checked = false;
                    }
                });

                updateSelectedGames(); // Обновляем текст кнопки
            })
            .catch(error => {
                console.error('Ошибка при обновлении списка игр:', error);
            });
    }
});

document.getElementById('tariff').addEventListener('change', function () {
    const tariffId = this.value;
    const dropdownList = document.getElementById('dropdown-general-games');
    const button = document.querySelector(".dropdown-button");

    if (!tariffId) {
        // Если тариф не выбран, показываем все игры
        Array.from(dropdownList.querySelectorAll('label')).forEach(label => {
            label.style.display = '';  // Показываем все игры
            const input = label.querySelector('input');
            input.disabled = false;
        });
        updateSelectedGames(); // Обновляем текст кнопки
        return;
    }

    fetch(`/admin/booking/edit/general-games?tariffId=${encodeURIComponent(tariffId)}`)
        .then(response => response.json())
        .then(selectedGames => {
            const selectedGameIds = new Set(selectedGames.map(game => game.Id));

            // Прячем или показываем игры в зависимости от того, есть ли их Id в ответе от сервера
            Array.from(dropdownList.querySelectorAll('label')).forEach(label => {
                const input = label.querySelector('input');
                const gameId = input.value;

                if (selectedGameIds.has(gameId)) {
                    // Если игра в ответе от сервера, показываем её
                    label.style.display = '';
                    input.disabled = false;
                } else {
                    // Если игры нет в ответе от сервера, скрываем её
                    label.style.display = 'none';
                    input.disabled = true;
                    input.checked = false;
                }
            });

            updateSelectedGames(); // Обновляем текст кнопки
        })
        .catch(error => {
            console.error('Ошибка при обновлении списка игр:', error);
        });
});

// Функция для отображения количества выбранных игр
function updateSelectedGames() {
    const selectedGames = document.querySelectorAll("input[name='general-games']:checked");
    const button = document.querySelector(".dropdown-button");

    // Обновляем текст кнопки в зависимости от количества выбранных игр
    const selectedCount = selectedGames.length;
    if (selectedCount === 0) {
        button.textContent = "Выберите общие игры (необязательно)";
    } else {
        button.textContent = `Выбрано игр: ${selectedCount}`;
    }
}

// Функция для переключения отображения выпадающего списка
function toggleDropdown() {
    const dropdown = document.getElementById("dropdown-general-games");
    dropdown.style.display = dropdown.style.display === "block" ? "none" : "block";
}

// Закрытие выпадающего окна при клике вне его
window.addEventListener("click", function (event) {
    const dropdown = document.getElementById("dropdown-general-games");
    const button = document.querySelector(".dropdown-button");
    if (!button.contains(event.target) && !dropdown.contains(event.target)) {
        dropdown.style.display = "none";
    }
});

// Обработчик изменения состояния флажков
const checkboxes = document.querySelectorAll("input[name='general-games']");
checkboxes.forEach(checkbox => {
    checkbox.addEventListener("change", updateSelectedGames);
});
