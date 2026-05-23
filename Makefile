.PHONY: run build tidy

run:
	go run ./cmd/main.go

build:
	go build -o server ./cmd/main.go

tidy:
	go mod tidy
