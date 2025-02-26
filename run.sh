#!/bin/bash

check_and_update_go_version() {

    current_version=$(go version | awk '{print $4}' | sed 's/go//g')
    if [ $? -ne 0 ]; then
        echo "Error: Failed to get the current Go version."
        return 1
    fi

    latest_version=$(curl -s -f https://golang.org/VERSION?m=text)
    if [ $? -ne 0 ]; then
        echo "Error: Failed to fetch the latest Go version."
        return 1
    fi

    echo "Current Go version: $current_version"
    echo "Latest Go version: $latest_version"

    if [ "$current_version" != "$latest_version" ]; then
        echo "Updating go.mod to use the latest version: $latest_version"
        go mod edit -go="$latest_version"
        if [ $? -ne 0 ]; then
            echo "Error: Failed to update the Go version in go.mod."
            return 1
        fi
        go mod tidy
    else
        echo "You are using the latest Go version: $current_version"
    fi
}

if [ -f 'go.mod' ]; then
    go mod tidy
    check_and_update_go_version
    xdg-open http://localhost:33333 || xdg-open --new-tab http://localhost:33333 ||  firefox http://localhost:33333
    go run main.go
    if [ $? -ne 0 ]; then
        echo "Application failed to run."
    fi
else
    echo "go.mod file not found"
    echo "installing go module..."
    go mod init forum
    if [ $? -ne 0 ]; then
        echo "Failed to initialize go module."
        exit 1
    fi
    go mod tidy
    check_and_update_go_version
    go run main.go
    if [ $? -ne 0 ]; then
        echo "Application failed to run."
    fi
fi
