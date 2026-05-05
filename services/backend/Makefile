APP=owlchat-backend

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP) ./cmd/server

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal
