SHELL := /bin/bash

default: docker-compose

.PHONY: install-tools
install-tools:
	go install github.com/cespare/reflex@latest

.PHONY: reset-local-storage
reset-local-storage:
	$(eval LOCAL_DB_PATH := $(shell source .env && echo $$LOCAL_DB_PATH))
	$(eval LOCAL_DB_NAME := $(shell source .env && echo $$LOCAL_DB_NAME))
	$(eval LOCAL_FILES_PATH := $(shell source .env && echo $$LOCAL_FILES_PATH))
	-rm -rf local
	mkdir -p $(LOCAL_DB_PATH)
	mkdir -p $(LOCAL_FILES_PATH)
	sqlite3 $(LOCAL_DB_PATH)/$(LOCAL_DB_NAME) < migrations/init.sql

.PHONY: bin/files-api
bin/files-api:
	go build -o bin/files-api cmd/files-api/main.go

.PHONY: start
start: bin/files-api
	source .env && bin/files-api

.PHONY: dev
dev: install-tools
	reflex -s -R 'local/|bin/' make start

.PHONY: test
test:
	GIN_MODE=test go test -v ./...

.PHONY: test
docker:
	docker build -t files-api .

.PHONY: docker-compose
docker-compose:
	source .env && docker-compose up