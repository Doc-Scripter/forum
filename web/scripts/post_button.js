//=====  create post container event listener which enables  scrolling ==========
document.addEventListener("DOMContentLoaded", function () {
    const createPostBtn = document.querySelector(".create-post-btn");
    const footer = document.querySelector("footer");
    function checkButtonVisibility() {
        const footerRect = footer.getBoundingClientRect();
        const buttonRect = createPostBtn.getBoundingClientRect();
        if (buttonRect.bottom > footerRect.top) {
            createPostBtn.classList.add("hidden");
        } else {
            createPostBtn.classList.remove("hidden");
        }
    }
    window.addEventListener("scroll", checkButtonVisibility);
    window.addEventListener("resize", checkButtonVisibility);
});
