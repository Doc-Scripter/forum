// DOM Elements
const postsContainer = document.querySelector(".posts-container");
const filterBtns = document.querySelectorAll(".filter-btn");
const categoryFilter = document.getElementById("category-filter");
const menuContent = document.querySelector(".menu-content");
const createPostBtn = document.querySelector(".create-post-btn");
const modal = document.querySelector(".modal");
const closeModal = document.querySelector(".close-modal");
const postForm = document.querySelector(".post-form");
const hamburgerIcon = document.querySelector(".hamburger-icon");
const likebutton = document.querySelectorAll(".like-btn");
const dislikebutton = document.querySelectorAll(".dislike-btn");

let currentFilter = "allPosts";
let currentCategory = "all";
let posts = [];
let route = "/posts";

postsContainer.addEventListener("click", (e) => {
  if (e.target.classList.contains("like-btn")) {
    console.log("like button clicked");
    const postId = e.target.dataset.postId;
    fetch("/likes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id: postId }),
    })
      // .then(response => response.json())
      .then(() => {
        fetchPosts(route);
      })
      .catch(
        (error) => console.error(error)
        // alert('You have already liked this post')
      );
  }

  if (e.target.classList.contains("dislike-btn")) {
    console.log("dislike button clicked");

    const postId = e.target.dataset.postId;
    fetch("/dislikes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id: postId }), // Ensure userId is defined
    })
      // .then(response => response.json())
      // .then(data => console.log(data))
      .then(() => {
        fetchPosts(route);
      })
      .catch(
        (error) => console.error(error)
        // alert('You have already disliked this post')
      );
  }
});

postForm.addEventListener("submit", (e) => {
  alert("Post submitted successfully!");
  // e.preventDefault();
  // fetchPosts();
  // postForm.reset()
});

// State

function fetchPosts(route) {
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
  if (filteredPosts === null) {
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
    // <span class="post-date">${post.date}</span>
    // <span class="post-author">By ${post.author}</span>
    // <div class="post-footer">
    // <div class="post-actions">
    //        <button class="action-btn like-btn ${post.liked ? 'active' : ''}" data-id="${post.id}">
    //         👍 ${post.likes}
    //       </button>
    //        <button class="action-btn dislike-btn ${post.disliked ? 'active' : ''}" data-id="${post.id}">
    //        👎 ${post.dislikes}
    //        </button>
    //        <button class="comments-toggle" data-post-id="${post.id}">
    //     💬 Comments (${post.comments.length})
    //   </button>
    //   </div>
    postsContainer.innerHTML = filteredPosts
      .map(
        (post) => `
      <article class="post">
      <div class="post-header">
      </div>
      <h2 class="post-category">${post.category}</h2>
      
       <h2 class="post-title">${post.title}</h2>
      <p class="post-content">${post.content}</p>
    <div class="post-footer">
     <div class="post-actions">
            <button class="action-btn like-btn" data-post-id=${post.post_id}>
             👍${post.likes}
           </button>
            <button class="action-btn dislike-btn" data-post-id=${post.post_id}>
            👎${post.dislikes}
            </button>
            
       </div>
    </article>
    `
      )
      .join("");
  }
}

document.addEventListener("DOMContentLoaded", function () {
  const usernameElement = document.getElementById("username");
  const username = usernameElement.textContent; // Use textContent to get the text

  function getInitials(name) {
    let initials = name
      .split(" ")
      .map((word) => word.charAt(0).toUpperCase())
      .join("");
    return initials;
  }

  const initials = getInitials(username);
  document.getElementById("profile-Name").textContent = initials; // Update the avatar
});

function toggleSidebar() {
  const sidebar = document.getElementById("sidebar");
  sidebar.classList.toggle("open");
}

document.addEventListener("DOMContentLoaded", fetchPosts("/posts"));

// Create Post Modal
createPostBtn.addEventListener("click", () => {
  modal.classList.add("active");
});

closeModal.addEventListener("click", () => {
  modal.classList.remove("active");
});

// modal.addEventListener('click', (e) => {
//   if (e.target === modal) {
//     modal.classList.remove('active');
//   }
// });

//   posts.unshift(newPost);
//   modal.classList.remove('active');
//   postForm.reset();
//   filterPosts();
// });

// <button class="comments-toggle" data-post-id="${post.id}">
//   💬 Comments (${post.comments.length})
// </button>
//     <div class="comments-section" id="comments-${post.id}">
//       ${post.comments.map(comment => `
//         <div class="comment">
//           <strong>${comment.author}</strong>
//           <p>${comment.content}</p>
//           <small>${comment.date}</small>
//         </div>
//       `).join('')}
//       <form class="comment-form" data-post-id="${post.id}">
//         <input type="text" class="comment-input" placeholder="Add a comment..." required>
//         <button type="submit" class="comment-submit">Comment</button>
//       </form>
//     </div>
//   </article>
// `).join('');

// // Add comments toggle functionality
// document.querySelectorAll('.comments-toggle').forEach(btn => {
//   btn.addEventListener('click', (e) => {
//     const postId = btn.dataset.postId;
//     const commentsSection = document.getElementById(`comments-${postId}`);
//     commentsSection.classList.toggle('active');
//   });
// });
// }

// function handleLike(btn) {
//   const postId = parseInt(btn.dataset.id);
//   const post = posts.find(p => p.id === postId);
//   if (post) {
//     if (post.disliked) {
//       post.disliked = false;
//       post.dislikes--;
//     }
//     post.liked = !post.liked;
//     post.likes += post.liked ? 1 : -1;
//     filterPosts();
//   }
// }

// function handleDislike(btn) {
//   const postId = parseInt(btn.dataset.id);
//   const post = posts.find(p => p.id === postId);
//   if (post) {
//     if (post.liked) {
//       post.liked = false;
//       post.likes--;
//     }
//     post.disliked = !post.disliked;
//     post.dislikes += post.disliked ? 1 : -1;
//     filterPosts();
//   }
// }

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
