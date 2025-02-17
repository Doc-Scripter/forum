// Function to close notification
function closeNotification() {
    const notification = document.getElementById('notification');
    notification.classList.remove('show');
}

// Function to show notification
function showNotification(message, type = 'error') {
    const notification = document.getElementById('notification');
    const messageSpan = document.getElementById('notification-message');
    
    if (!notification || !messageSpan) {
        console.error('Notification elements not found');
        return;
    }

    messageSpan.textContent = message;
    notification.className = `notification ${type}`;
    
    // Show the notification
    notification.classList.add('show');

    // Hide the notification after 3 seconds
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

// Client-side validation
function validateForm() {
    const email = document.querySelector('input[name="email"]').value;
    const password = document.querySelector('input[name="password"]').value;

    if (!email || !password) {
        showNotification('Email and password are required');
        return false;
    }

    if (!email.includes('@')) {
        showNotification('Please enter a valid email address');
        return false;
    }

    if (password.length < 6) {
        showNotification('Password must be at least 6 characters');
        return false;
    }

    return true;
}

// Form submission handler
document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            if (!validateForm()) {
                return;
            }

            try {
                const formData = new FormData(e.target);
                const response = await fetch('/logging', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams(formData)
                });
                console.log(formData);

                const errorMessage = await response.text();

                if (!response.ok) {
                    throw new Error(errorMessage || 'Login failed');
                }

                // If successful, redirect to home
                window.location.href = '/home';
            } catch (error) {
                showNotification(error.message);
                // Clear password field on error
                document.querySelector('input[name="password"]').value = '';
            }
        });
    }

    // Check for error parameter in URL and show notification if present
    const urlParams = new URLSearchParams(window.location.search);
    const errorMsg = urlParams.get('error');
    if (errorMsg) {
        showNotification(decodeURIComponent(errorMsg));
    }
});

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
