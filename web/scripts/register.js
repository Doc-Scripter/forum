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

function validatePassword(password) {
    // Check minimum length
    if (password.length < 6) {
        showNotification('Password must be at least 6 characters long');
        return false;
    }

    // Check for uppercase letter
    if (!/[A-Z]/.test(password)) {
        showNotification('Password must contain at least one uppercase letter');
        return false;
    }

    // Check for lowercase letter
    if (!/[a-z]/.test(password)) {
        showNotification('Password must contain at least one lowercase letter');
        return false;
    }

    // Check for number
    if (!/[0-9]/.test(password)) {
        showNotification('Password must contain at least one number');
        return false;
    }

    // Check for special character
    if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
        showNotification('Password must contain at least one special character');
        return false;
    }

    return true;
}

// Client-side validation
function validateForm() {
    const username = document.querySelector('input[name="username"]').value;
    const email = document.querySelector('input[name="email"]').value;
    const password = document.querySelector('input[name="password"]').value;
    const confirmPassword = document.querySelector('input[name="confirm_password"]').value;

    if (!username || !email || !password || !confirmPassword) {
        showNotification('All fields are required');
        return false;
    }

    if (username.length < 3) {
        showNotification('Username must be at least 3 characters long');
        return false;
    }

    if (!email.includes('@')) {
        showNotification('Please enter a valid email address');
        return false;
    }

    // Check password requirements
    if (!validatePassword(password)) {
        return false;
    }

    if (password !== confirmPassword) {
        showNotification('Passwords do not match');
        return false;
    }

    return true;
}

// Form submission handler
document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            if (!validateForm()) {
                return;
            }

            try {
                const formData = new FormData(e.target);
                const response = await fetch('/register', {
                    method: 'POST',
                    body: formData
                });

                const errorMessage = await response.text();

                if (!response.ok) {
                    throw new Error(errorMessage || 'Registration failed');
                }

                // If successful, show success message and redirect
                showNotification('Registration successful! Redirecting to login...', 'success');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } catch (error) {
                showNotification(error.message);
                // Clear password fields on error
                document.querySelector('input[name="password"]').value = '';
                document.querySelector('input[name="confirm_password"]').value = '';
            }
        });
    }

    // Check for error parameter in URL
    const urlParams = new URLSearchParams(window.location.search);
    const errorMsg = urlParams.get('error');
    if (errorMsg) {
        showNotification(decodeURIComponent(errorMsg));
    }
}); 