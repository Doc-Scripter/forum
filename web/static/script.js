// The slide in menu
document.addEventListener("DOMContentLoaded", function () {
  const sidebar = document.getElementById("sidebar");

  // Toggle sidebar
  window.toggleSidebar = function() {
    sidebar.classList.toggle("active");
  };

  // Set profile picture initials dynamically
  const username = "John Smith"; // Replace with dynamic user data
  const initials = username.split(" ").map(n => n[0]).join("").toUpperCase();
  document.getElementById("profile-pic").innerText = initials;
  document.getElementById("username").innerText = username;
});
