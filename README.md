- [forum](#forum)
- [Overview](#overview)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Docker](#docker)
- [License](#license)


# FORUM

## Overview
This project involves creating a web forum that facilitates user communication through posts and comments. The forum supports features like associating categories with posts, liking and disliking posts and comments, and filtering posts based on various criteria. The project utilizes SQLite for database management and Docker for containerization.

## Features
* User Authentication: Users can register and log in. Sessions are managed using cookies with expiration dates.

* Post and Comment Management: Registered users can create posts and comments. Posts can be associated with one or more categories.

* Likes and Dislikes: Registered users can like or dislike posts and comments. The counts are visible to all users.

* Filtering: Users can filter posts by categories, created posts, and liked posts.

## Technologies Used

* Backend: Go (Golang)

* Database: SQLite

* Containerization: Docker

* Authentication: UUID (Bonus), bcrypt for password encryption (Bonus)

## Getting Started
## Requirements
* Go 1.17 or higher
* SQLite 3.36 or higher
* Docker 20.10 or higher

## Installation
1. Clone the repository: 
```bash
git clone https://learn.zone01kisumu.ke/git/shaokoth/forum
cd forum
```

2. Build and run the docker image
- In the terminal
```bash
./script.sh
```
Open your web browser and navigate to 
```bash
http://localhost:8080
```

# Alternatively 
1. Run the program directly through the terminal

```bash
go run .
```

2. Access the Forum
Open your web browser and navigate to 
```bash
http://localhost:33333
```
## Database Schema
### Tables
* Users
    - id (Primary Key)
    - uuid (Unique)
    - username (Unique)
    - email (Unique)
    - password (Encrypted)
* Posts
    - post_id (Primary Key)
    - user_uuid (Foreign Key to Users)
    - title
    - filename
    - content
    - filepath
    - comments
    - category
    - created_at
* Comments
    - id (Primary Key)
    - post_id (Foreign Key to Posts)
    - user_uuid (Foreign Key to Users)
    - content
    - likes
    - dislikes
    - created_at
* LikesDislikes
    - id (Primary Key)
    - post_id (Foreign Key to Posts)
    - comment_id (Foreign Key to comments)
    - user_uuid (Foreign Key to Users)
    - like_dislike
* Sessions
    - id (Primary Key)
    - user_id (Foreign Key to users(id))
    - session_token (Unique)
    - expires_at

## API Endpoints
## Authentication
* POST/register
    - Registers a new user.
* POST/login
    - Logs in a user and creates a session.
* POST /logout
    - logs out a user and deletes a session.

## Posts
* GET /posts
    - Retrieves all posts.

* POST /create-post
    - Creates a new post.
* GET /favorites
    - Displays liked posts

## Comments
* POST /addcomment
    - Adds a comment to a post.
* GET /comments
    - Display all comments for a specific post.

## LikesDislikes
* POST /likes
    - Likes a post.
* POST /dislikes
    - Dislikes a post.
* POST /likesComment
    - Likes a comment.
* POST /dislikeComment
    - Dislikes a comment.

## Error Handling
The application handles various HTTP status errors and technical errors gracefully. Custom error messages are provided for better debugging and user feedback.

# Testing
Unit tests are written using the Go testing framework and can be run using the go test command.
```bash
go test ./...
```
## or
```bash
make test 
``` 

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your changes.

# License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
