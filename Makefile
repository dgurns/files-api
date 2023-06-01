default: docker-compose

.PHONY: install-tools
install-tools:
	go install github.com/cespare/reflex@latest

bin/files-api:
	go build -o bin/files-api

.PHONY: env
env:
	$(source .env)
	@echo "Environment variables set from .env file"

dev: install-tools bin/files-api env
	reflex -R 'gen/|bin/' -s bin/files-api

docker:
	docker build -t files-api .

docker-compose: env
	docker-compose build
	docker-compose up