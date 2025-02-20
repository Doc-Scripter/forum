(function () {
  // Push a new state to the history
  window.history.pushState({}, "", window.location.href);

  // Add an event listener for the popstate event
  window.addEventListener("popstate", function () {
    // Redirect to the login page or show a session expired message
    window.location.href = "/";
  });
})();
