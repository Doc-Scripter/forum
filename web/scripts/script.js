// DOM Elements
const postsContainer = document.querySelector('.posts-container');
const filterBtns = document.querySelectorAll('.filter-btn');
const categoryFilter = document.getElementById('category-filter');
// const userInitials = document.querySelector('.user-initials');
const menuContent = document.querySelector('.menu-content');
// const logoutBtn = document.querySelector('.logout-btn');
const createPostBtn = document.querySelector('.create-post-btn');
const modal = document.querySelector('.modal');
const closeModal = document.querySelector('.close-modal');
const postForm = document.querySelector('.post-form');
const hamburgerIcon = document.querySelector('.hamburger-icon');
const likebutton= document.querySelectorAll('.like-btn')
const dislikebutton= document.querySelectorAll('.dislike-btn')


// Event listeners
postsContainer.addEventListener('click', (e) => {
  if (e.target.classList.contains('like-btn')) {
    console.log('like button clicked');
    const postId = e.target.dataset.postId;
    fetch("/likes", {
      method: "POST",
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ post_id: postId }),
    })
      // .then(response => response.json())
      .then(() => {
        fetchPosts()
      })
      .catch(error => 
        console.error(error)
        // alert('You have already liked this post')
      );
     
      
  }
  
  if (e.target.classList.contains('dislike-btn')) {
    console.log('dislike button clicked');

    const postId = e.target.dataset.postId;
    fetch("/dislikes", {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ post_id: postId}), // Ensure userId is defined
    })
    // .then(response => response.json())
    // .then(data => console.log(data))
    .then(()=>{
      fetchPosts()
    })
    .catch(error => 
      console.error(error)
      // alert('You have already disliked this post')
    );
  }
});

function loadComments(postId) {
  fetch(`/comments?post_id=${postId}`)
    .then(response => response.json())
    .then(comments => {
      const commentList = document.querySelector(`#comments-${postId} .comments-list`);
      commentList.innerHTML = comments.length
        ? comments.map(comment => `<p><strong>${comment.author}:</strong> ${comment.content}</p>`).join('')
        : "<p>No comments yet. Be the first to comment!</p>";
    })
    .catch(error => console.error("Error loading comments:", error));
}

// Event listener for comment button (displaying the comment section)
postsContainer.addEventListener("click", (e) => {
  if (e.target.classList.contains("comment-btn")) {
    const postId = e.target.dataset.postId;
    const commentSection = document.getElementById(`comments-${postId}`);
    
    // Toggle the comment section
    if (commentSection.style.display === "none") {
      commentSection.style.display = "block";
      fetchComments(postId);  // Fetch comments when shown
    } else {
      commentSection.style.display = "none";  // Hide comment section
    }
  }
});


// Submit comment form
postsContainer.addEventListener("submit", (e) => {
  if (e.target.classList.contains("comment-form")) {
    e.preventDefault();

    const postId = e.target.closest(".post").dataset.postId;
    const commentInput = e.target.querySelector(".comment-input");
    const commentText = commentInput.value.trim();

    if (commentText === "") return;

    fetch("/add-comments", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ post_id: postId, text: commentText }),
    })
    .then(() => {
      commentInput.value = "";
      fetchComments(postId);
    })
    .catch(error => console.error("Error adding comment:", error));
  }
});


postForm.addEventListener('submit', (e) => {
  alert('Post submitted successfully!');
})

// State
let currentFilter = 'all';
let currentCategory = 'all';
let posts = [];


  function fetchPosts(){
  fetch("/posts")
    .then((response) => response.json())
    .then((data) => {
       posts = data;
      displayPosts(posts,currentCategory);
    })
    .catch(error => { console.error("Error fetching posts:",error)});

}

// Fetch and display comments for a post
function fetchComments(postId) {
  fetch(`/comments?post_id=${postId}`)
    .then(response => response.json())
    .then(comments => {
      const commentsList = document.querySelector(`#comments-${postId} .comments-list`);
      commentsList.innerHTML = comments.length
        ? comments.map(comment => `<p><strong>${comment.author}:</strong> ${comment.text}</p>`).join('')
        : "<p>No comments yet. Be the first to comment!</p>";
    })
    .catch(error => console.error("Error fetching comments:", error));
}

// Display posts
function displayPosts(posts, category) {
  let filteredPosts = [];
  if (category === "all") {
    filteredPosts = posts;
  } else {
    filteredPosts = posts.filter(post => post.category === category);
  }

  if (!filteredPosts || filteredPosts.length === 0) {
    postsContainer.innerHTML = `
      <article class="post">
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
    postsContainer.innerHTML = filteredPosts.map(post => `
      <article class="post" data-post-id="${post.post_id}">
        <div class="post-header">
        </div>
        <h2 class="post-category">${post.category}</h2>
        <h2 class="post-title">${post.title}</h2>
        <p class="post-content">${post.content}</p>
        <div class="post-footer">
          <div class="post-actions">
            <button class="action-btn like-btn" data-post-id="${post.post_id}">
              👍 ${post.likes}
            </button>
            <button class="action-btn dislike-btn" data-post-id="${post.post_id}">
              👎 ${post.dislikes}
            </button>
            <button class="action-btn comment-btn" data-post-id="${post.post_id}">
              💬 ${post.comments_count || 0}
            </button>
          </div>
          <div class="comments-section" id="comments-${post.post_id}" style="display: none;">
            <div class="comments-list"></div>
            <form class="comment-form">
              <input type="text" class="comment-input" placeholder="Write a comment..." required>
              <button type="submit" class="submit-comment-btn">Post</button>
            </form>
          </div>
        </div>
      </article>
    `).join('');
  }
}



document.addEventListener("DOMContentLoaded", fetchPosts);

// Close menu when clicking outside
document.addEventListener('click', (e) => {
  if (!e.target.closest('.hamburger-menu')) {
    menuContent.classList.remove('active');
    hamburgerIcon.classList.remove('active');
  }else{
    menuContent.classList.add('active');
    hamburgerIcon.classList.add('active');
  }
});


// Create Post Modal
createPostBtn.addEventListener('click', () => {
  modal.classList.add('active');
});

closeModal.addEventListener('click', () => {
  modal.classList.remove('active');
});


// Filter posts
filterBtns.forEach(btn => {
  btn.addEventListener('click', () => {
    filterBtns.forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    currentFilter = btn.dataset.filter;
    displayPosts(posts,currentFilter);
  });
});

categoryFilter.addEventListener('change', (e) => {
  currentCategory = e.target.value;
  displayPosts(posts,currentCategory);
  // filterPosts();
});