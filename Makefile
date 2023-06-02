default: docker-compose

.PHONY: install-tools
install-tools:
	go install github.com/cespare/reflex@latest

bin/files-api:
	go build -o bin/files-api cmd/files-api/main.go

start:
	source .env && bin/files-api

dev: install-tools bin/files-api
	reflex -R 'gen/|bin/' -s make start

docker:
	docker build -t files-api .

docker-compose:
	docker-compose build
	source .env && docker-compose up