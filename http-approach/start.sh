#!/bin/bash

if ! go list -m github.com/google/uuid &>/dev/null; then
    echo "Installing github.com/google/uuid..."
    go get -u github.com/google/uuid
else
    echo "github.com/google/uuid is already installed, skipping..."
fi

if ! go list -m golang.org/x/crypto &>/dev/null; then
    echo "Installing golang.org/x/crypto..."
    go get -u golang.org/x/crypto
else
    echo "golang.org/x/crypto is already installed, skipping..."
fi

# Run the program
go run main.go
