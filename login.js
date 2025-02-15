// Function to show notificationfunction showNotification(message) {
function ShowNotification(message){
    console.log('we got this far');
    const notification = document.getElementById('notification');

    // Update the text content of the notification element
    notification.textContent = message;

    // Show the notification
    notification.classList.add('show');

    // Hide the notification after 3 seconds (3000 milliseconds)
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

// Client-side validation
function ValidateForm() {
    
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    if (!email || !password) {
        ShowNotification('Email and password are required');
        return false;
    }
    if (!email.includes('@')) {
        ShowNotification('Please enter a valid email address');
        return false;
    }
    if (password.length < 6) {
        ShowNotification('Password must be at least 6 characters');
        return false;
    }

    if (password !== "anxile#123") {
        ShowNotification('Invalid password');
        return false;
    }
    
    
    return true;
}

document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent default form submission
    if (ValidateForm()) {
        // Proceed with form submission, e.g., using AJAX or:
        this.submit(); // Submit the form if validation passes
    }
});
