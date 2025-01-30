#!/bin/bash

# Install dependencies
go mod init forum
go get github.com/gofiber/fiber/v2
go get gorm.io/driver/sqlite
go get gorm.io/gorm
go get github.com/gofiber/storage/sqlite3
go get github.com/gofiber/fiber/v2/middleware/session
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt

# start the project 
go run main.go