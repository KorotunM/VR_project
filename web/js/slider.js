document.addEventListener("DOMContentLoaded", () => {
    const sliderContainer = document.querySelector(".slider-container");
    const prevButton = document.querySelector(".slider-btn.prev");
    const nextButton = document.querySelector(".slider-btn.next");

    if (!sliderContainer || !prevButton || !nextButton) {
        console.error("Не найдены элементы слайдера. Проверь HTML-разметку.");
        return;
    }

    const cardWidth = document.querySelector(".tariff-card").offsetWidth + 20; // ширина + gap
    const totalCards = sliderContainer.children.length;
    const visibleCards = 3; // Количество видимых карточек
    let currentIndex = 0; // Текущий индекс

    // Изначально задаём CSS-свойство transition для анимации
    sliderContainer.style.transition = "transform 0.3s ease-in-out";

    // Кнопка "вправо"
    nextButton.addEventListener("click", () => {
        if (currentIndex < totalCards - visibleCards) {
            currentIndex++;
            updateSlider();
        }
        toggleButtons();
    });

    // Кнопка "влево"
    prevButton.addEventListener("click", () => {
        if (currentIndex > 0) {
            currentIndex--;
            updateSlider();
        }
        toggleButtons();
    });

    // Обновление позиции слайдера
    function updateSlider() {
        const offset = -currentIndex * cardWidth;
        sliderContainer.style.transform = `translateX(${offset}px)`;
    }

    // Функция для скрытия кнопок на границах
    function toggleButtons() {
        prevButton.style.visibility = currentIndex === 0 ? "hidden" : "visible";
        nextButton.style.visibility = currentIndex >= totalCards - visibleCards ? "hidden" : "visible";
    }

    // Изначально скрываем кнопку "влево", если слайдер в начале
    toggleButtons();
});

function scrollToBooking() {
    const bookingForm = document.querySelector(".booking-form");
    if (bookingForm) {
        bookingForm.scrollIntoView({ behavior: "smooth" });
    }
}
