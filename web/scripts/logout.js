(function () {
  window.history.pushState({}, "", window.location.href);

  window.addEventListener("popstate", function () {
    window.location.href = "/";
  });
})();
