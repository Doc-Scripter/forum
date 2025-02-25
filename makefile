test:
	go test -cover  ./...
format:
	gofmt -w -s .
run :
	go run .
docker:
	mkdir -p ~/forum_database
	docker build -t myapp .
	docker run -p 8080:33333 -v ~/forum_database:/root/forum.db myapp