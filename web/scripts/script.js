// DOM Elements
const postsContainer = document.querySelector('.posts-container');
// const filterBtns = document.querySelectorAll('.filter-btn');
// const categoryFilter = document.getElementById('category-filter');
// const userInitials = document.querySelector('.user-initials');
// const menuContent = document.querySelector('.menu-content');
// const logoutBtn = document.querySelector('.logout-btn');
const createPostBtn = document.querySelector('.create-post-btn');
const modal = document.querySelector('.modal');
const closeModal = document.querySelector('.close-modal');
const postForm = document.querySelector('.post-form');
// const hamburgerIcon = document.querySelector('.hamburger-icon');


// State
// let currentFilter = 'all';
// let currentCategory = 'all';
// let posts = [];

// filterbutton.addEventListener("click", () => {
  // const category = document.getElementById("category-filter").value;

  // const filterCriteria = {
  //   category,
  // };
  // Fetch posts from API, then update the UI with the filtered posts
  function fetchPosts(){
  fetch("/posts")
    .then((response) => response.json())
    .then((data) => {
      const posts = data;
      displayPosts(posts);
      return posts;
      // const filteredPosts = filterPosts(posts, categoryFilter);
    })
    .catch(error => { console.error("Error fetching posts:",error)});
  // postsContainer.innerHTML = "";
  // filterPosts.forEach((post) => {
  //   const postElement = document.createElement("div");
  //   postElement.innerHTML = "<h2>${post.title}</h2><p>${post.content}</p>";
  //   postsContainer.appendChild(postElement);
  // });
}

function displayPosts(posts) {
  console.log(posts);
  if (posts===null||posts.length===0||posts==[]){
    postsContainer.innerHTML = `<article class="post">
    <div class="post-header">
    <span class="post-date">NO Date</span>
    </div>
    <h2 class="post-title">No posts available</h2>
    <p class="post-content">No posts to display</p>
     <div class="post-footer">
        <span class="post-author"></span>
      </div>
    </article>`
    ;
  }else{
    
    // <span class="post-date">${post.date}</span>
    // <h2 class="post-title">${post.title}</h2>
    // <span class="post-author">By ${post.author}</span>
  postsContainer.innerHTML = posts.map(post=>`
    <article class="post">
    <div class="post-header">
    </div>
    <h2 class="post-category">${post.category}</h2>

    <p class="post-content">${post.content}</p>
     <div class="post-footer">
      </div>
    </article>
    `).join('')
  }
}


document.addEventListener("DOMContentLoaded", fetchPosts);



// });

// Initialize user information
// document.querySelector('.user-name').textContent = currentUser.name;

// Toggle hamburger menu
// /

// Close menu when clicking outside
// document.addEventListener('click', (e) => {
//   if (!e.target.closest('.hamburger-menu')) {
//     menuContent.classList.remove('active');
//     hamburgerIcon.classList.remove('active');
//   }
// });

// Logout functionality
// logoutBtn.addEventListener('click', () => {
//   alert('Logged out successfully!');
  // Add actual logout logic here
// });

// Filter posts
// function filterPosts() {
//   let filteredPosts = posts;

//   // Apply category filter
//   if (currentCategory !== 'all') {
//     filteredPosts = filteredPosts.filter(post => post.category === currentCategory);
//   }

//   // Apply post type filter
//   if (currentFilter === 'created') {
//     filteredPosts = filteredPosts.filter(post => post.authorId === currentUser.id);
//   } else if (currentFilter === 'liked') {
//     filteredPosts = filteredPosts.filter(post => post.liked);
//   }

//   displayPosts(filteredPosts);
// }

// function filterPosts(posts, filter) {
  // const filteredPosts = [];

  // posts.forEach((post) => {
  //   if (
  //     (filter.category === undefined || post.category === filter.category) &&
  //     (filter.dateRange === undefined ||
  //       (post.date >= filter.dateRange.start &&post.date <= filter.dateRange.end)) &&
  //     (filter.likesRange === undefined ||(post.likes >= filter.likesRange.min &&
  //         post.likes <= filter.likesRange.max))
  //   ) {
  //     filteredPosts.push(post);
  //   }
  // });

  // return filteredPosts;
// }
//??

// Create Post Modal
createPostBtn.addEventListener('click', () => {
  modal.classList.add('active');
});

closeModal.addEventListener('click', () => {
  modal.classList.remove('active');
});

// modal.addEventListener('click', (e) => {
//   if (e.target === modal) {
//     modal.classList.remove('active');
//   }
// });

// Handle post creation
postForm.addEventListener('submit', (e) => {
//   e.preventDefault();
//   const formData = new FormData(postForm);
//   const newPost = {
//     id: posts.length + 1,
//     title: formData.get('title'),
//     content: formData.get('content'),
//     category: formData.get('category'),
//     author: currentUser.name,
//     authorId: currentUser.id,
//     likes: 0,
//     dislikes: 0,
//     liked: false,
//     disliked: false,
//     date: new Date().toISOString().split('T')[0],
//     comments: []
displayPosts();
});
  
//   posts.unshift(newPost);
//   modal.classList.remove('active');
//   postForm.reset();
//   filterPosts();
// });

// Display posts
// function displayPosts(postsToDisplay) {
  // postsContainer.innerHTML = postsToDisplay.map(post => `
  //   <article class="post">
  //     <div class="post-header">
  //       <span class="post-date">${post.date}</span>
  //     </div>
  //     <h2 class="post-title">${post.title}</h2>
  //     <p class="post-content">${post.content}</p>
  //     <div class="post-footer">
  //       <span class="post-author">By ${post.author}</span>
  //       <div class="post-actions">
  //         <button class="action-btn like-btn ${post.liked ? 'active' : ''}" data-id="${post.id}">
  //           👍 ${post.likes}
  //         </button>
  //         <button class="action-btn dislike-btn ${post.disliked ? 'active' : ''}" data-id="${post.id}">
  //           👎 ${post.dislikes}
  //         </button>
  //       </div>
  //     </div>
  //     <button class="comments-toggle" data-post-id="${post.id}">
  //       💬 Comments (${post.comments.length})
  //     </button>
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

  // // Add interaction handlers
  // document.querySelectorAll('.like-btn').forEach(btn => {
  //   btn.addEventListener('click', () => handleLike(btn));
  // });

  // document.querySelectorAll('.dislike-btn').forEach(btn => {
  //   btn.addEventListener('click', () => handleDislike(btn));
  // });

  // document.querySelectorAll('.comment-form').forEach(form => {
  //   form.addEventListener('submit', handleComment);
  // });

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

// function handleComment(e) {
//   e.preventDefault();
//   const postId = parseInt(e.target.dataset.postId);
//   const post = posts.find(p => p.id === postId);
//   const input = e.target.querySelector('.comment-input');
  
//   if (post && input.value.trim()) {
//     post.comments.push({
//       id: post.comments.length + 1,
//       author: currentUser.name,
//       content: input.value.trim(),
//       date: new Date().toISOString().split('T')[0]
//     });
//     filterPosts();
//     input.value = '';
//   }
// }

// Event listeners for filters
// filterBtns.forEach(btn => {
//   btn.addEventListener('click', () => {
//     filterBtns.forEach(b => b.classList.remove('active'));
//     btn.classList.add('active');
//     currentFilter = btn.dataset.filter;
//     filterPosts();
//   });
// });

// categoryFilter.addEventListener('change', (e) => {
//   currentCategory = e.target.value;
//   filterPosts();
// });

// Initial display
// filterPosts();