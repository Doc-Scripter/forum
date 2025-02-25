package handlers

import (
	"encoding/json"
	d "forum/database"
	m "forum/models"
	"sort"
	"fmt"
	"net/http"
	e "forum/Error"
)

// ==== This function will handle post creation and insertion of the post into the database ====

func PostsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	
	rows, err := d.Db.Query("SELECT title,content,created_at,post_id,filename,filepath FROM posts")
	
	if err != nil {
		ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return

	}

	defer rows.Close()

	var posts []m.Post
	for rows.Next() {

		var eachPost m.Post

		err := rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_id, &eachPost.Filename, &eachPost.Filepath)
		eachPost.Seperate_Categories()
		if err != nil {
			ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		defer rows.Close()

		eachPost.CommentsCount = commentsCount

		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		var comments []m.Comment
		for rows.Next() {
			var comment m.Comment
			rows.Scan(&comment.Content)
			comments = append(comments, comment)
		}

		eachPost.Comments = comments

		var likeCount, dislikeCount int

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", &eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", &eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}

	postsJson, err := json.Marshal(OrderPosts(posts))
	if err != nil {
		ErrorPage(fmt.Errorf("|post handler| ---> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return
	}

	e.LOGGER("[SUCCESS]: Fetching posts was a success!", nil)
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}


func OrderPosts(posts []m.Post) []m.Post{
	sort.Slice(posts, func(i, j int) bool {
        return posts[i].CreatedAt.After(posts[j].CreatedAt)
    })
    return posts
}