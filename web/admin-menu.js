document.addEventListener('DOMContentLoaded', () => {
    const toggleMenuButton = document.getElementById('toggleMenu');
    const sidebar = document.getElementById('sidebar');

    toggleMenuButton.addEventListener('click', () => {
        sidebar.classList.toggle('hidden');
    });
});