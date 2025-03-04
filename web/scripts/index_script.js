//====== Declare the DOM elements that will be required =========
const postsContainer = document.querySelector(".posts-container");
const filterBtns = document.querySelectorAll(".filter-btn");
const categoryFilter = document.getElementById("category-filter");
const menuContent = document.querySelector(".menu-content");
const comments = document.querySelectorAll(".comments-section");

let currentFilter = "allPosts";
let currentCategory = "all";
let posts = [];
let route = "/posts";

// ===== function to fetch the posts from the backend endpoint ======
function fetchPosts(route) {
  console.log("I got here");
  fetch(route)
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      displayPosts(data, currentCategory);
    })
    .catch((error) => {
      console.error("Error fetching posts:", error);
    });
}

//============The function that splits the string coming from the backend and displays them in different labels of a post============
function createCategoryElements(categories) {
  if (!Array.isArray(categories) || categories.length === 0) {
    return "";
  }

  let html = '<div class="category-container">';

  categories.forEach(category => {
    const trimmedCategory = category.trim();

    if (trimmedCategory.length > 0) {
      html += `<h2 class="post-category">${trimmedCategory}</h2>`;
    }
  });

  html += '</div>';

  return html;
}

//===== function to display the posts  and their respective comments =====
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
            💬 ${
              post.comments?.length === 1
                ? `${post.comments.length} Comment`
                : `${post.comments?.length || 0} Comments`
            }
          </button>

        </div>
        <div class="comments-section" id="comments-${post.post_id}">
        
        </div>

          <form class="comment-form">
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
  
  const commentForms = document.querySelectorAll('form');
  
commentForms.forEach((form) => {
    form.addEventListener("submit", (e) => {
      e.preventDefault(); // Prevent the default form submission
      if (confirm("Please Login to comment! Click OK to go to login page.")) {
        window.location.href = "/login"; // Redirect to login page
      }
    });
  });
  

 //=================Add comments toggle functionality===================
 document.querySelectorAll(".comments-toggle").forEach((btn) => {
  btn.addEventListener("click", (e) => {
    const postId = btn.dataset.postId;
    const commentsSection = document.getElementById(`comments-${postId}`);
    
    if (commentsSection.classList.contains("active")) {
      commentsSection.classList.remove("active");
      commentsSection.style.display = "none";
    } else {
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
      👍${comment.likes}
      </button>
      <button class="comment dislikeBtn" data-comment-id="${
        comment.comment_id
      }">
      👎${comment.dislikes}
      </button>
      </div>
      `
      )
      .join(``);
    attachCommentActionListeners(element);
  }
}


//============ comment actions ===============
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

//====== function to fetch the comments from the backend endpoint ======
function fetchComments(element, commentId) {
  fetch("/comments", {
    method: "POST",
    content: "application/json",
    body: JSON.stringify({ comment_id: commentId }),
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(  data);
      displayComments(data, element);
    })
    .catch((error) => {
      console.error("Error fetching comments:", error);
    });
}

//=======================Function to display the comments=======================
function displayComments(comments, element) {
  element.innerHTML = comments
    .map(
      (comment) => `
  <div class="comment"><p>${comment.content}</p></div>
    <div class="comment-actions">
    <button class="comment likeBtn" data-comment-id="${comment.comment_id}">
      👍${comment.likes}
      </button>
      <button class="comment dislikeBtn" data-comment-id="${comment.comment_id}">
      👎${comment.dislikes}
      </button>
    </div>
    `
    )
    .join(``);
  attachCommentActionListeners(element);
}

document.addEventListener("DOMContentLoaded", fetchPosts("/posts"));

//======= Event listeners for filter buttons ======
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

//===== Event listener to filter the displayed posts ======
categoryFilter.addEventListener("change", (e) => {
  currentCategory = e.target.value;
  fetchPosts(route);
});

//======= Theme Toggle Functionality =======
document.addEventListener("DOMContentLoaded", () => {
  const themeToggle = document.getElementById("theme-toggle");

  // ======= Check for saved theme preference, default to light if none  ========
  const savedTheme = localStorage.getItem("theme") || "light-theme";
  document.body.classList.toggle("light-theme", savedTheme === "light-theme");

  // Ensure light theme is applied by default
  if (!localStorage.getItem("theme")) {
    document.body.classList.add("light-theme");
    localStorage.setItem("theme", "light-theme");
  }

  //======= Update icon visibility based on current theme =======
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

//=========== Function to update theme icon =========
function updateThemeIcon() {

  const sunIcon = document.querySelector(".sun");
  const moonIcon = document.querySelector(".moon");
  const isLightTheme = document.body.classList.contains("light-theme");

  sunIcon.style.display = isLightTheme ? "none" : "block";
  moonIcon.style.display = isLightTheme ? "block" : "none";
}
