.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t go-microservice .

docker-run:
	docker run -p 8080:8080 go-microservice

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down