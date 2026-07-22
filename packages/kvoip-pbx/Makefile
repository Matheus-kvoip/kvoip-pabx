.PHONY: run build test tidy fmt vet

APP=kvoip-pbx
BIN=bin/$(APP)

run:
	go run ./cmd/kvoip-pbx

build:
	mkdir -p bin
	go build -o $(BIN) ./cmd/kvoip-pbx

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...
