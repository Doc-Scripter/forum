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
                    body: formData
                });

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
