// Sample data arrays for post generation
const topics = {
    'Science': [
        { title: 'Understanding Gravity Waves', content: 'Recent discoveries in gravitational wave detection have revolutionized our understanding of space-time...' },
        { title: 'The Physics of Black Holes', content: 'Exploring how gravity behaves at the event horizon of black holes...' },
        { title: 'Gravity and Space Travel', content: 'How gravitational assists help spacecraft navigate the solar system...' },
        { title: 'Einstein\'s Theory of Gravity', content: 'Breaking down general relativity and its implications for modern physics...' }
    ],
    'Math': [
        { title: 'Introduction to Combinatorics', content: 'Understanding the basics of counting and arrangement problems...' },
        { title: 'Permutations vs Combinations', content: 'Exploring the key differences and applications in real-world scenarios...' },
        { title: 'Graph Theory in Combinatorics', content: 'How graphs help solve complex counting problems...' },
        { title: 'Probability and Combinatorics', content: 'The intersection of probability theory and combinatorial mathematics...' }
    ],
    'Technology': [
        { title: 'AI Fundamentals Explained', content: 'A beginner\'s guide to understanding artificial intelligence...' },
        { title: 'Machine Learning Basics', content: 'Introduction to how machines learn from data...' },
        { title: 'Neural Networks Simplified', content: 'Understanding the building blocks of deep learning...' },
        { title: 'Future of AI Technology', content: 'Predictions and trends in artificial intelligence development...' }
    ],
    'Health': [
        { title: 'Understanding Mental Health', content: 'Breaking down the basics of mental wellness...' },
        { title: 'Nutrition Fundamentals', content: 'Essential nutrients for a healthy lifestyle...' }
    ],
    'Sports': [
        { title: 'Training Techniques', content: 'Advanced methods for athletic performance...' },
        { title: 'Sports Psychology', content: 'Mental preparation for competitive sports...' }
    ]
};

function generateRandomPosts(count = 30) {
    const posts = [];
    const categories = Object.keys(topics);
    
    for (let i = 0; i < count; i++) {
        // Select random category
        const category = categories[Math.floor(Math.random() * categories.length)];
        const topicList = topics[category];
        const topic = topicList[Math.floor(Math.random() * topicList.length)];
        
        // Generate random likes and dislikes
        const likes = Math.floor(Math.random() * 100);
        const dislikes = Math.floor(Math.random() * 50);
        
        // Generate random date within last 30 days
        const date = new Date();
        date.setDate(date.getDate() - Math.floor(Math.random() * 30));
        
        // Create post object
        const post = {
            post_id: i + 1,
            title: topic.title,
            content: topic.content,
            category: category.toLowerCase(),
            likes: likes,
            dislikes: dislikes,
            created_at: date.toISOString(),
            author: "User" + Math.floor(Math.random() * 100),
            userInteraction: {
                liked: false,
                disliked: false
            }
        };
        
        posts.push(post);
    }
    
    return posts;
}

// Example usage:
// const randomPosts = generateRandomPosts();
// console.log(JSON.stringify(randomPosts, null, 2));

// If you want to use this with your existing fetch function:
function mockFetchPosts() {
    return new Promise((resolve) => {
        const posts = generateRandomPosts();
        resolve(posts);
    });
}

// Add these new functions to handle likes and dislikes
function handleLike(post) {
    if (post.userInteraction.disliked) {
        // If post was disliked, remove dislike and add like
        post.dislikes--;
        post.likes++;
        post.userInteraction.disliked = false;
        post.userInteraction.liked = true;
    } else if (post.userInteraction.liked) {
        // If already liked, remove like
        post.likes--;
        post.userInteraction.liked = false;
    } else {
        // If neither liked nor disliked, add like
        post.likes++;
        post.userInteraction.liked = true;
    }
    return post;
}

function handleDislike(post) {
    if (post.userInteraction.liked) {
        // If post was liked, remove like and add dislike
        post.likes--;
        post.dislikes++;
        post.userInteraction.liked = false;
        post.userInteraction.disliked = true;
    } else if (post.userInteraction.disliked) {
        // If already disliked, remove dislike
        post.dislikes--;
        post.userInteraction.disliked = false;
    } else {
        // If neither liked nor disliked, add dislike
        post.dislikes++;
        post.userInteraction.disliked = true;
    }
    return post;
}

// Example usage:
let post = posts[0];
post = handleLike(post);    // Adds like
post = handleLike(post);    // Removes like (toggle off)
post = handleDislike(post); // Adds dislike

// Export for use in other files
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { generateRandomPosts, mockFetchPosts, handleLike, handleDislike };
}

