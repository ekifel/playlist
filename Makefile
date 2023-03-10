.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/playlist/main.go

run:
	go run ./cmd/playlist

lint:
	golangci-lint run

test:
	go test --short -v ./...
	make test.coverage

docker-build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down
