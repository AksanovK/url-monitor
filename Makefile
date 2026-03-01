.PHONY: run test lint build docker

run:
	go run cmd/server/main.go

test:
	go test ./... -v -race

lint:
	golangci-lint run

build:
	go build -o bin/server cmd/server/main.go

docker-build:
	docker build -t url-monitor .

docker-run:
	docker run -p 8080:8080 url-monitor