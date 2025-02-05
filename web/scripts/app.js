const { createApp, ref, computed } = Vue

createApp({
    setup() {
        const categories = ['Technology', 'Lifestyle', 'Gaming', 'Science', 'Art']
        
        const posts = ref([
            { 
                id: 1, 
                userId: 1, 
                category: 'Technology', 
                content: 'AI is amazing!', 
                likes: 0,
                comments: []
            },
            { 
                id: 2, 
                userId: 2, 
                category: 'Gaming', 
                content: 'New game release!', 
                likes: 0,
                comments: []
            }
        ])

        const currentUser = ref({ 
            id: 1, 
            name: 'John Doe', 
            email: 'john.doe@example.com',
            avatar: 'https://via.placeholder.com/100' 
        })

        const currentView = ref('posts')
        const selectedCategory = ref('')
        const showCreatePostModal = ref(false)
        const isSidebarOpen = ref(false)

        const newPost = ref({
            category: '',
            content: ''
        })

        const newComment = ref({
            postId: null,
            content: ''
        })

        function toggleSidebar() {
            isSidebarOpen.value = !isSidebarOpen.value
        }

        function handleSidebarAction(action) {
            switch(action) {
                case 'logout':
                    // Implement logout logic
                    alert('Logout clicked')
                    break
            }
            // Close sidebar after action
            isSidebarOpen.value = false
        }

        const filteredPosts = computed(() => {
            return posts.value.filter(post => 
                !selectedCategory.value || post.category === selectedCategory.value
            )
        })

        const userPosts = computed(() => 
            posts.value.filter(post => post.userId === currentUser.value.id)
        )

        const likedPosts = computed(() => 
            posts.value.filter(post => post.likes > 0)
        )

        function createPost() {
            if (!newPost.value.category || !newPost.value.content) {
                alert('Please fill all fields')
                return
            }

            posts.value.push({
                id: posts.value.length + 1,
                userId: currentUser.value.id,
                category: newPost.value.category,
                content: newPost.value.content,
                likes: 0,
                comments: []
            })

            // Reset modal
            newPost.value.category = ''
            newPost.value.content = ''
            showCreatePostModal.value = false
        }

        function likePost(postId) {
            const post = posts.value.find(p => p.id === postId)
            if (post) {
                post.likes++
            }
        }

        function addComment(postId, commentContent) {
            const post = posts.value.find(p => p.id === postId)
            if (post) {
                post.comments.push({
                    id: post.comments.length + 1,
                    userId: currentUser.value.id,
                    userName: currentUser.value.name,
                    content: commentContent,
                    timestamp: new Date()
                })
            }
        }

        return {
            categories,
            posts,
            currentUser,
            currentView,
            selectedCategory,
            showCreatePostModal,
            newPost,
            newComment,
            filteredPosts,
            userPosts,
            likedPosts,
            createPost,
            likePost,
            addComment,
            isSidebarOpen,
            toggleSidebar,
            handleSidebarAction
        }
    },
    components: {
        'post-list': {
            props: ['posts'],
            data() {
                return {
                    commentInputs: {}
                }
            },
            methods: {
                initCommentInput(postId) {
                    if (!this.commentInputs[postId]) {
                        this.$set(this.commentInputs, postId, '')
                    }
                },
                submitComment(postId) {
                    const commentContent = this.commentInputs[postId]
                    if (commentContent && commentContent.trim()) {
                        this.$emit('add-comment', postId, commentContent)
                        this.commentInputs[postId] = ''
                    }
                }
            },
            template: `
                <div class="post-list">
                    <div v-for="post in posts" :key="post.id" class="post">
                        <div class="post-header">
                            <span class="category">{{ post.category }}</span>
                        </div>
                        <div class="post-content">
                            {{ post.content }}
                        </div>
                        <div class="post-footer">
                            <button @click="$emit('like-post', post.id)">
                                👍 {{ post.likes }}
                            </button>
                        </div>
                        
                        <!-- Comments Section -->
                        <div class="comments-section">
                            <h4>Comments</h4>
                            <div v-if="post.comments.length === 0" class="no-comments">
                                No comments yet
                            </div>
                            <div v-for="comment in post.comments" :key="comment.id" class="comment">
                                <strong>{{ comment.userName }}</strong>
                                <p>{{ comment.content }}</p>
                            </div>
                            
                            <!-- Comment Input -->
                            <div class="comment-input">
                                <textarea 
                                    v-model="commentInputs[post.id]" 
                                    @focus="initCommentInput(post.id)"
                                    placeholder="Write a comment..."
                                ></textarea>
                                <button @click="submitComment(post.id)">
                                    Post Comment
                                </button>
                            </div>
                        </div>
                    </div>
                    <p v-if="posts.length === 0">No posts found</p>
                </div>
            `
        }
    }
}).mount('#app')