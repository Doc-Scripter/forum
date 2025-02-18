// DOM Elements
const postsContainer = document.querySelector(".posts-container");
const filterBtns = document.querySelectorAll(".filter-btn");
const categoryFilter = document.getElementById("category-filter");
const menuContent = document.querySelector(".menu-content");
const comments = document.querySelectorAll(".comments-section");

let currentFilter = "allPosts";
let currentCategory = "all";
let posts = [];
let route = "/posts";





// State
function fetchPosts(route) {
  console.log("here", route);
  fetch(route)
    .then((response) => response.json())
    .then((data) => {
      displayPosts(data, currentCategory);
    })
    .catch((error) => {
      console.error("Error fetching posts:", error);
    });
}

function displayPosts(posts, category) {
  let filteredPosts = [];
  if (category === "all") {
    filteredPosts = posts;
  } else {
    if (posts !== null) {
      filteredPosts = posts.filter((post) => post.category === category);
    }
  }
  if (filteredPosts === null || !posts || filteredPosts.length === 0) {
    postsContainer.innerHTML = `<article class="post">
    <div class="post-header">
    <span class="post-date">NO Date</span>
    </div>
    <h2 class="post-title">No posts available</h2>
    <p class="post-content">No posts to display</p>
     <div class="post-footer">
        <span class="post-author"></span>
      </div>
    </article>`;
  } else {
    postsContainer.innerHTML = filteredPosts
      .map(
        (post) => `
      <article class="post">
      <div class="post-header"></div>
      <h2 class="post-category">${post.category}</h2>
      <h2 class="post-title">${post.title}</h2>
      <p class="post-content">${post.content}</p>
      <div class="post-footer">
        <div class="post-actions">
          <button class="action-btn like-btn" data-post-id="${post.post_id}">
            👍${post.likes}
          </button>
          <button class="action-btn dislike-btn" data-post-id="${post.post_id}">
            👎${post.dislikes}
          </button>
          <button class="comments-toggle" data-post-id="${post.post_id}">
            💬 ${post.comments?.length === 1 ? `${post.comments.length} Comment` : `${post.comments?.length || 0} Comments`}
          </button>

        </div>
        <div class="comments-section" id="comments-${post.post_id}">
          ${
            post.comments
              ? post.comments
                  .map(
                    (comment) => `
          <div class="comment">
            <p>${comment}</p>
          </div>
          `
                  )
                  .join("")
              : ""
          }

          <form
            class="comment-form"
            data-post-id="${post.post_id}"
            action="/register"
            method="get"
          >
            <input type="hidden" name="post_id" value="${post.post_id}" />
            <input
              type="text"
              name="add-comment"
              class="comment-input"
              placeholder="Add a comment..."
              required
            />
            <button type="submit" class="comment-submit">Comment</button>
          </form>
        </div>
      </div>
    </article>
       `
      )
      .join("");
  }

  // Add comments toggle functionality
  document.querySelectorAll(".comments-toggle").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      const postId = btn.dataset.postId;
      const commentsSection = document.getElementById(`comments-${postId}`);
      commentsSection.classList.toggle("active");
      fetch("/comments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ post_id: postId }),
      }).catch((error) => console.error(error));
    });
  });
}

document.addEventListener("DOMContentLoaded", fetchPosts("/posts"));


// Event listeners for filters
filterBtns.forEach((btn) => {
  btn.addEventListener("click", () => {
    filterBtns.forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    currentFilter = btn.dataset.filter;
    switch (currentFilter) {
      case "allPosts":
        route = "/posts";
        break;
      case "created":
        route = "/myPosts";
        break;
      case "liked":
        route = "/favorites";
        break;
      default:
        console.error("Invalid filter value");
    }
    console.log("new route", route);
    fetchPosts(route);
  });
});

categoryFilter.addEventListener("change", (e) => {
  currentCategory = e.target.value;
  console.log("category", currentCategory);
  fetchPosts(route);
});

// Theme Toggle Functionality
document.addEventListener("DOMContentLoaded", () => {
  const themeToggle = document.getElementById("theme-toggle");

  // Check for saved theme preference, default to light if none
  const savedTheme = localStorage.getItem("theme") || "light-theme";
  document.body.classList.toggle("light-theme", savedTheme === "light-theme");

  // Ensure light theme is applied by default
  if (!localStorage.getItem("theme")) {
    document.body.classList.add("light-theme");
    localStorage.setItem("theme", "light-theme");
  }

  // Update icon visibility based on current theme
  updateThemeIcon();

  themeToggle.addEventListener("click", () => {
    document.body.classList.toggle("light-theme");
    updateThemeIcon();

    // Save theme preference
    const currentTheme = document.body.classList.contains("light-theme")
      ? "light-theme"
      : "dark-theme";
    localStorage.setItem("theme", currentTheme);
  });
});

// Function to update theme icon
function updateThemeIcon() {
  const sunIcon = document.querySelector(".sun");
  const moonIcon = document.querySelector(".moon");
  const isLightTheme = document.body.classList.contains("light-theme");

  // Show moon in light theme (to switch to dark)
  // Show sun in dark theme (to switch to light)
  sunIcon.style.display = isLightTheme ? "none" : "block";
  moonIcon.style.display = isLightTheme ? "block" : "none";
}

