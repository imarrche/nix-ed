.PHONY: build
build:
	go build -o ./build/api ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

start:
	./build/api
