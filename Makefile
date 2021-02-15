.PHONY: build
build:
	go build -o ./build/api ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

start:
	./build/api

test:
	go test -cover -coverprofile=coverage.html -timeout 30s ./...

.PHONY: coverage
coverage:
	go tool cover -html=coverage.html
