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