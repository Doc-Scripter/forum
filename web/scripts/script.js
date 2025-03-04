//====== Declare the DOM elements that will be required =========
const postsContainer = document.querySelector(".posts-container");
const filterBtns = document.querySelectorAll(".filter-btn");
const categoryFilter = document.getElementById("category-filter");
const createPostBtn = document.querySelector(".create-post-btn");
const modal = document.querySelector(".modal");
const closeModal = document.querySelector(".close-modal");
const postForm = document.querySelector(".post-form");
const comments = document.querySelector(".comments-section");

const commentform = document.querySelectorAll(".comment-form");

let currentFilter = "allPosts";
let currentCategory = "all";
let posts = [];
let route = "/posts";

// ===== function to process a post request =========
postsContainer.addEventListener("click", (e) => {
  if (e.target.classList.contains("like-btn")) {
    const postId = e.target.dataset.postId;
    const isCurrentlyLiked = e.target.classList.contains("active");

    fetch("/likes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id: postId }),
    })
      .then((response) => {
        if (!response.ok) {
          return response.text().then((text) => {
            throw new Error(text || "Failed to like post");
          });
        }
        //====== Remove active class from dislike button =================
        const dislikeBtn = e.target.parentElement.querySelector(".dislike-btn");
        dislikeBtn.classList.remove("active");

        //====== Toggle active class on like button based on current state =========
        if (isCurrentlyLiked) {
          e.target.classList.remove("active");
        } else {
          e.target.classList.add("active");
        }
        return fetchPosts(route);
      })
      .catch((error) => {
        showNotification(error.message, "error");
        console.error(error);
      });
  }

  if (e.target.classList.contains("dislike-btn")) {
    const postId = e.target.dataset.postId;
    const isCurrentlyDisliked = e.target.classList.contains("active");

    fetch("/dislikes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id: postId }),
    })
      .then((response) => {
        if (!response.ok) {
          return response.text().then((text) => {
            throw new Error(text || "Failed to dislike post");
          });
        }
        //====== Remove active class from like button =================
        const likeBtn = e.target.parentElement.querySelector(".like-btn");
        likeBtn.classList.remove("active");

        //====== Toggle active class on dislike button based on current state =================
        if (isCurrentlyDisliked) {
          e.target.classList.remove("active");
        } else {
          e.target.classList.add("active");
        }
        return fetchPosts(route);
      })
      .catch((error) => {
        showNotification(error.message, "error");
        console.error(error);
      });
  }
});

postForm.addEventListener("submit", (e) => {
  e.preventDefault();

  //====== Get form values =================
  const title = postForm.querySelector('[name="title"]').value.trim();
  const content = postForm.querySelector('[name="content"]').value.trim();
  const categories = Array.from(
    postForm.querySelectorAll('[name="category"]:checked')
  ).map((cb) => cb.value);

  //====== Validation checks =================
  if (!title) {
    showNotification("Post title is required", "error");
    return;
  }

  if (!content) {
    showNotification("Post content is required", "error");
    return;
  }

  if (categories.length === 0) {
    showNotification("Please select at least one category", "error");
    return;
  }

  //====== If validation passes, proceed with form submission ===========
  const formData = new FormData(postForm);

  fetch("/create-post", {
    method: "POST",
    body: formData,
  })
    .then((response) => {
      if (!response.ok) {
        showNotification("Failed to create post", "error");
        return;
      }
      showNotification("Post submitted successfully!", "success");
      modal.classList.remove("active");
      postForm.reset();
      return fetchPosts(route);
    })
    .catch((error) => {
      showNotification("Failed to create post", "error");
      console.error("Error creating post:", error);
    });
});

//========== This is the function that will be fetching posts ==== =================
function fetchComments(element, commentId) {
  fetch("/comments", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ comment_id: commentId }),
  })
    .then((response) => response.json())
    .then((data) => {
      displayComments(data, element);
    })
    .catch((error) => {
      console.error("Error fetching posts:", error);
    });
}

//===== Function to create HTML for categories from an array ====
function createCategoryElements(categories) {
  if (!Array.isArray(categories) || categories.length === 0) {
    return "";
  }

  let html = '<div class="category-container">';

  categories.forEach((category) => {
    const trimmedCategory = category.trim();

    if (trimmedCategory.length > 0) {
      html += `<h2 class="post-category">${trimmedCategory}</h2>`;
    }
  });

  html += "</div>";

  return html;
}

//====== Function to fetch posts from the backend =====
function fetchPosts(route) {
  console.log(route);
  fetch(route)
    .then((response) => response.json())
    .then((data) => {
      displayPosts(data, currentCategory);
    })
    .catch((error) => {
      console.error("Error fetching posts:", error);
    });
}

//===== This function will display the posts ========
function displayPosts(posts, category) {
  let filteredPosts = [];
  if (category === "all") {
    filteredPosts = posts;
  } else {
    if (posts !== null) {
      filteredPosts = posts.filter((post) => post.category.includes(category));
    }
  }
  if (filteredPosts === null || !posts || filteredPosts.length === 0) {
    postsContainer.innerHTML = `<article class="post">
    <div class="post-header">
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
        <div> ${createCategoryElements(post.category)} </div>
      <h2 class="post-title">${escapeHTML(post.title)}</h2>
      <p class="post-content">${escapeHTML(post.content)}</p>
      ${
        post.filepath
          ? `<div class="post-image-container">
          <img src="/image/${post.filepath}" alt="${post.filename}" class="post-image">
          </div>`
          : ``
      }

      <div class="post-footer">
        <div class="post-actions">
          <button class="action-btn like-btn" data-post-id="${post.post_id}">
            👍${post.likes}
          </button>
          <button class="action-btn dislike-btn" data-post-id="${post.post_id}">
            👎${post.dislikes}
          </button>
          <button class="comments-toggle" data-post-id="${post.post_id}">
            💬 ${
              post.comments?.length === 1
                ? `${post.comments.length} Comment`
                : `${post.comments?.length || 0} Comments`
            }
          </button>
          </div>
        </div>
        <div class="comments-section" id="comments-${post.post_id}">
        </div>

          <form
            class="comment-form"
            data-post-id="${post.post_id}"
            action="/addcomment"
            method="post"
             onsubmit="return validateComment(event)"
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
    </article>
      `
      )
      .join("");
  }

  //=================Add comments toggle functionality===================
  document.querySelectorAll(".comments-toggle").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      const postId = btn.dataset.postId;
      const commentsSection = document.getElementById(`comments-${postId}`);

      // Toggle the active class
      if (commentsSection.classList.contains("active")) {
        // If comments are showing, hide them
        commentsSection.classList.remove("active");
        commentsSection.style.display = "none";
      } else {
        // If comments are hidden, show them and fetch
        commentsSection.classList.add("active");
        commentsSection.style.display = "block";

        fetch("/comments", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ post_id: postId }),
        })
          .then((res) => res.json())
          .then((data) => {
            displayComments(data, commentsSection);
          })
          .catch((error) => console.error(error));
      }
    });
  });
}

/* ====== Escapes special characters in a string to prevent XSS attacks when rendering HTML.//+
 * @param {string} str - The input string to be escaped.//+
 * @returns {string} The escaped string with special characters replaced by their HTML entity equivalents. ==== */
function escapeHTML(str) {
  return str.replace(/[&<>"'/]/g, function (char) {
    switch (char) {
      case "&":
        return "&amp;";
      case "<":
        return "&lt;";
      case ">":
        return "&gt;";
      case '"':
        return "&quot;";
      // case "'":
      //   return "&#039;";
      case "/":
        return "&#47;";
      default:
        return char;
    }
  });
}

//========= Function to display the comments =========
function displayComments(comments, element) {
  console.log("This is the comments: ", comments);
  if (comments && comments !== null) {
    element.innerHTML = comments
      .map(
        (comment) => `
      <div class="comment"><p>${escapeHTML(comment.content)}</p></div>
      <div class="comment-actions">
      <button class="comment likeBtn" data-comment-id="${comment.comment_id}">
      👍 ${comment.likes}
      </button>
      <button class="comment dislikeBtn" data-comment-id="${
        comment.comment_id
      }">
      👎 ${comment.dislikes}
      </button>
      </div>
      `
      )
      .join(``);
    attachCommentActionListeners(element);
  }
}

//==========================comment actions====================
function attachCommentActionListeners(container) {
  container.querySelectorAll(".comment-actions").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      if (e.target.classList.contains("likeBtn")) {
        const commentId = e.target.dataset.commentId;
        fetch("/likesComment", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ comment_Id: commentId }),
        })
          .then(() => fetchComments(container, commentId))
          .catch((error) => console.error(error));
      }
      if (e.target.classList.contains("dislikeBtn")) {
        // e.stopPropagation();
        const commentId = e.target.dataset.commentId;
        fetch("/dislikesComment", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ comment_id: commentId }),
        })
          .then(() => fetchComments(container, commentId))
          .catch((error) => console.error(error));
      }
    });
  });
}

document.addEventListener("DOMContentLoaded", fetchPosts("/posts"));

// ======== Create Post Modal =========
createPostBtn.addEventListener("click", () => {
  modal.classList.add("active");
});
closeModal.addEventListener("click", () => {
  modal.classList.remove("active");
});
modal.addEventListener("click", (e) => {
  if (e.target === modal) {
    modal.classList.remove("active");
  }
});

//========== Event listeners for filters =========
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
    fetchPosts(route);
  });
});
categoryFilter.addEventListener("change", (e) => {
  currentCategory = e.target.value;
  fetchPosts(route);
});

//===================Theme Toggle Functionality====================
document.addEventListener("DOMContentLoaded", () => {
  const themeToggle = document.getElementById("theme-toggle");
  const savedTheme = localStorage.getItem("theme") || "light-theme";
  document.body.classList.toggle("light-theme", savedTheme === "light-theme");
  if (!localStorage.getItem("theme")) {
    document.body.classList.add("light-theme");
    localStorage.setItem("theme", "light-theme");
  }
  updateThemeIcon();
  themeToggle.addEventListener("click", () => {
    document.body.classList.toggle("light-theme");
    updateThemeIcon();
    const currentTheme = document.body.classList.contains("light-theme")
      ? "light-theme"
      : "dark-theme";
    localStorage.setItem("theme", currentTheme);
  });
});

//======= Function to update theme icon =========
function updateThemeIcon() {
  const sunIcon = document.querySelector(".sun");
  const moonIcon = document.querySelector(".moon");
  const isLightTheme = document.body.classList.contains("light-theme");
  sunIcon.style.display = isLightTheme ? "none" : "block";
  moonIcon.style.display = isLightTheme ? "block" : "none";
}

//====== For the checkboxes ==============
const checkboxGroup = document.querySelector(".checkbox-group");
const myDivs = document.querySelectorAll(".category-option");

myDivs.forEach((div) => {
  div.style.backgroundColor = "#ccc";
  div.style.color = "#000";
});

checkboxGroup.addEventListener("change", function () {
  const checkedCount = [...checkboxGroup.querySelectorAll("input:checked")]
    .length;
  myDivs.forEach((div) => {
    if (div.querySelector("input").checked) {
      div.classList.add("checked");
      div.style.backgroundColor = "green";
      div.style.color = "#fff";
    } else {
      // if (checkedCount >= 3) {
      //     div.style.backgroundColor = "#f00"; // Red
      //     div.style.color = "#fff";
      // } else {
      div.style.backgroundColor = "#ccc";
      div.style.color = "#000";
      // }
    }
  });
});

//===== function at the top with other utility functions ========
function showNotification(message, type = "success") {
  const notification = document.getElementById("notification");
  notification.textContent = message;
  notification.className = `notification ${type}`;

  setTimeout(() => {
    notification.classList.add("show");
  }, 100);

  setTimeout(() => {
    notification.classList.remove("show");
  }, 3000);
}

//========== Add validation for comment forms ========
document.querySelectorAll(".comment-form").forEach((form) => {
  form.addEventListener("submit", (e) => {
    const commentInput = form.querySelector(".comment-input");

    if (commentInput.value.trim() === "") {
      e.preventDefault();
      alert("Comment cannot be empty");
      return;
    }

    commentInput.value = commentInput.value.trim();
  });
});

//====== function to validate comment =====
function validateComment(event) {
  const form = event.target;
  const commentInput = form.querySelector(".comment-input");

  if (commentInput.value.trim() === "") {
    event.preventDefault();
    alert("Comment cannot be empty or contain only whitespace");
    commentInput.value = "";
    commentInput.focus();
    return false;
  }

  commentInput.value = commentInput.value.trim();
  return true;
}
