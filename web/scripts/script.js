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



// likebutton.forEach(btn => {
//   btn.addEventListener('click', () => {
//     console.log('like button clicked');
//   const postId = btn.dataset.postId;
//   fetch("/likes", {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ post_id: postId}),
//   })
//     .then(response => response.json())
//     .then(data => console.log(data))
//     .catch(error => console.error(error));
// });
// })

// dislikebutton.forEach(btn => {
//   btn.addEventListener('click', () => {
//     const postId = btn.dataset.postId;
//     fetch("/dislikes", {
//       method: 'POST',
//       headers: { 'Content-Type': 'application/json' },
//       body: JSON.stringify({ post_id: postId, user_id: userId }),
//     })
//       .then(response => response.json())
//       .then(data => console.log(data))
//       .catch(error => console.error(error));
//   });
//   })

  // const postsContainer = document.getElementById('postsContainer');

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
      .then(data => console.log(data))
      .catch(error => console.error(error));
  }
  
  if (e.target.classList.contains('dislike-btn')) {
    const postId = e.target.dataset.postId;
    fetch("/dislikes", {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ post_id: postId, user_id: userId }), // Ensure userId is defined
    })
      .then(response => response.json())
      .then(data => console.log(data))
      .catch(error => console.error(error));
  }
});

postForm.addEventListener('submit', (e) => {
  alert('Post submitted successfully!');
  // e.preventDefault();
  // fetchPosts();
  // postForm.reset()
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

function displayPosts(posts,category) {
  let filteredPosts = [];
  if (category==="all"){
    filteredPosts=posts
  }else{
    filteredPosts=posts.filter(post=>post.category===category);

  }
  
  if (filteredPosts===null||!posts){
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
    postsContainer.innerHTML = filteredPosts.map(post=>`
      <article class="post">
      <div class="post-header">
      </div>
      <h2 class="post-category">${post.category}</h2>
      
       <h2 class="post-title">${post.title}</h2>
      <p class="post-content">${post.content}</p>
    <div class="post-footer">
     <div class="post-actions">
            <button class="action-btn like-btn" data-post-id=${post.id}>
             👍 ${post.likes}
           </button>
            <button class="action-btn dislike-btn" data-post-id=${post.id}>
            👎 ${post.dislikes}
            </button>
            
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
document.addEventListener('click', (e) => {
  if (!e.target.closest('.hamburger-menu')) {
    menuContent.classList.remove('active');
    hamburgerIcon.classList.remove('active');
  }else{
    menuContent.classList.add('active');
    hamburgerIcon.classList.add('active');
  }
});

// Logout functionality
// logoutBtn.addEventListener('click', () => {
//   alert('Logged out successfully!');
  // Add actual logout logic here
// });

// Filter posts
// function filterPosts() {
//   let filteredPosts = posts;

//   // Apply category filter
//   // if (currentCategory !== 'all') {
//   //   filteredPosts = filteredPosts.filter(post => post.category === currentCategory);
//   // }
//   if (currentFilter === 'all') {
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
//   const filteredPosts = [];

//   posts.forEach((post) => {
//     if (
//       (filter.category === undefined || post.category === filter.category) &&
//       (filter.dateRange === undefined ||
//         (post.date >= filter.dateRange.start &&post.date <= filter.dateRange.end)) &&
//       (filter.likesRange === undefined ||(post.likes >= filter.likesRange.min &&
//           post.likes <= filter.likesRange.max))
//     ) {
//       filteredPosts.push(post);
//     }
//   });

//   return filteredPosts;
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

// Initial display
// filterPosts();