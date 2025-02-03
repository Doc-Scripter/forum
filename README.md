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


# forum

# Overview
This project is a web forum that allows users to communicate with each other, associate categories to posts, like and dislike posts and comments, and filter posts. The project uses SQLite as the database management system and is built using Go.

# Features
User registration and login
Post creation and commenting
Category association with posts
Liking and disliking posts and comments
Filtering posts by categories, created posts, and liked posts
Error handling for website errors and technical errors
Use of Docker for containerization

# Requirements
Go 1.17 or higher
SQLite 3.36 or higher
Docker 20.10 or higher
# Installation
Clone the repository: git clone https://github.com/your-username/forum-project.git
Change into the project directory: cd forum-project
Build the Docker image: docker build -t forum-project .
Run the Docker container: docker run -p 8080:8080 forum-project
# Usage
Open a web browser and navigate to http://localhost:8080
Register as a new user by filling out the registration form
Login to the forum using your registered credentials
Create a new post by filling out the post form
Associate categories with your post by selecting from the available categories
Like or dislike posts and comments by clicking on the like or dislike button
Filter posts by categories, created posts, or liked posts using the filter dropdown menu

# API Endpoints
/register: Register a new user
/login: Login to the forum
/posts: Create a new post
/posts/{id}: Get a post by ID
/posts/{id}/comments: Get comments for a post
/posts/{id}/like: Like a post
/posts/{id}/dislike: Dislike a post
/categories: Get all categories
/categories/{id}: Get a category by ID
/filter: Filter posts by categories, created posts, or liked posts

# Database Schema
The database schema is defined in the schema.sql file and consists of the following tables:

users: Stores user information
posts: Stores post information
comments: Stores comment information
categories: Stores category information
post_categories: Stores the association between posts and categories
likes: Stores like information for posts and comments
dislikes: Stores dislike information for posts and comments

# Testing
Unit tests are written using the Go testing framework and can be run using the go test command.

# Docker
The project uses Docker for containerization and can be built and run using the docker build and docker run commands.

# License
This project is licensed under the MIT License.