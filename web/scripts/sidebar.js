// Function to toggle body scroll
function toggleBodyScroll(disable) {
    document.body.style.overflow = disable ? 'hidden' : 'auto';
    
    // Also disable scrolling on the main content
    const postsContainer = document.querySelector('.posts-container');
    if (postsContainer) {
        postsContainer.style.overflow = disable ? 'hidden' : 'auto';
    }
    
    // Disable scrolling on any scrollable elements
    const scrollableElements = document.querySelectorAll('.post, .comments-section, .modal-content');
    scrollableElements.forEach(element => {
        element.style.overflow = disable ? 'hidden' : 'auto';
    });
}

// Update the sidebar toggle functionality
function toggleSidebar() {
    const sidebar = document.querySelector('.sidebar');
    if (sidebar) {
        const isOpen = sidebar.classList.toggle('open');
        toggleBodyScroll(isOpen); // Disable scroll when sidebar is open
    }
}

// Ensure scroll is re-enabled when clicking outside sidebar or pressing escape
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        const sidebar = document.querySelector('.sidebar');
        if (sidebar && sidebar.classList.contains('open')) {
            sidebar.classList.remove('open');
            toggleBodyScroll(false); // Re-enable scroll
        }
    }
});

// Optional: Click outside sidebar to close
document.addEventListener('click', function(event) {
    const sidebar = document.querySelector('.sidebar');
    const userIcon = document.querySelector('.user-icon'); // Assuming this is your sidebar toggle button
    
    if (sidebar && sidebar.classList.contains('open') && 
        !sidebar.contains(event.target) && 
        !userIcon.contains(event.target)) {
        sidebar.classList.remove('open');
        toggleBodyScroll(false); // Re-enable scroll
    }
});
