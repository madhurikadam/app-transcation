include .env.local
export

.PHONY: install
install:
	go mod download

# ----- Building -----
SERVICE ?= app-transcation

.PHONY: clean
clean:
	rm -f $(SERVICE)

# builds our binary
build:
	cd $(CURDIR); CGO_ENABLED=0 go build -ldflags '-w' -o .build/$(SERVICE) ./cmd

$(SERVICE): build

.PHONY: all
all: clean $(LINTER) build

docker-image:
	make build
	docker build -f Dockerfile  --rm -t $(SERVICE):local . --build-arg SERVICE=$(SERVICE)

dev:
	go run ./cmd/main.go

test:
	go test ./...

up:
	docker-compose up -d postgres

down:
	docker-compose down
