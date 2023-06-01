default: start

dev:
	go run cmd/files-api/main.go

bin/files-api:
	go build -o bin/files-api

docker:
	docker build -t files-api .

start:
	docker-compose up