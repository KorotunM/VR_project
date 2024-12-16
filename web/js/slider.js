document.addEventListener("DOMContentLoaded", () => {
    const sliderContainer = document.querySelector(".slider-container");
    const prevButton = document.querySelector(".slider-btn.prev");
    const nextButton = document.querySelector(".slider-btn.next");

    let currentIndex = 0;
    const cardWidth = document.querySelector(".tariff-card").offsetWidth + 20; // 20px - gap

    prevButton.addEventListener("click", () => {
        if (currentIndex > 0) {
            currentIndex--;
            updateSlider();
        }
    });

    nextButton.addEventListener("click", () => {
        if ((currentIndex + 1) * cardWidth < sliderContainer.scrollWidth) {
            currentIndex++;
            updateSlider();
        }
    });

    function updateSlider() {
        const offset = -currentIndex * cardWidth;
        sliderContainer.style.transform = `translateX(${offset}px)`;
    }
});

function scrollToBooking() {
    const bookingForm = document.querySelector(".booking-form");
    if (bookingForm) {
        bookingForm.scrollIntoView({ behavior: "smooth" });
    }
}
