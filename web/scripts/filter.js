const filterbutton = document.getElementById("filter-btn");
const postsContainer = document.getElementById("posts-container");

filterbutton.addEventListener("click", () => {
  const category = document.getElementById("category-filter").value;
  const startDate = document.getElementById("start-date").value;
  const endDate = document.getElementById("end-date").value;
  const minLikes = document.getElementById("min-likes").value;
  const maxLikes = document.getElementById("max-likes").value;

  const filterCriteria = {
    category,
    dateRange: { startDate: startDate, endDate: endDate },
    LikesRange: { min: minLikes, max: maxLikes },
  };
  // Fetch posts from API, then update the UI with the filtered posts
  fetch("/posts")
    .then((response) => response.json())
    .then((data) => {
      const posts = data;
      const filteredPosts = filterPosts(posts, filterCriteria);
    });

  postsContainer.innerHTML = "";
  filterPosts.forEach((post) => {
    const postElement = document.createElement("div");
    postElement.innerHTML = "<h2>${post.title}</h2><p>${post.content}</p>";
    postsContainer.appendChild(postElement);
  });
});

function filterPosts(posts, filter) {
  const filteredPosts = [];

  posts.forEach((post) => {
    if (
      (filter.category === undefined || post.category === filter.category) &&
      (filter.dateRange === undefined ||
        (post.date >= filter.dateRange.start &&post.date <= filter.dateRange.end)) &&
      (filter.likesRange === undefined ||(post.likes >= filter.likesRange.min &&
          post.likes <= filter.likesRange.max))
    ) {
      filteredPosts.push(post);
    }
  });

  return filteredPosts;
}
