# files-api

## Overview

This is a simple REST API for uploading, querying, and managing files. It
persists files locally via SQLite and the file system, but can be expanded to
support any cloud database and file storage provider.

## Requirements

- Go 1.20
  - [Download](https://golang.org/dl/)
- SQLite 3
  - Mac: `brew install sqlite3`
  - Linux: `sudo apt-get install sqlite3`
  - Windows: [Download](https://www.sqlite.org/download.html)
- Docker
  - [Download](https://www.docker.com/get-started)

## Running Locally

```sh
# Create an env file with defaults
cp .env.sample .env

# Set up your local SQLite and file storage
make reset-local-storage

# Run the API via docker-compose
make

# Make a curl request to upload a file
curl -X POST -H "Content-Type: multipart/form-data" \
	-F "file=@/path/to/file.pdf" \
	-F "metadata={\"key\":\"value\"}" \
	-u "demo:password" \
	http://localhost:8080/files/upload
```

## Endpoints

- See `openapi.yaml` for the full API spec

## Local Development

```sh
# Install tools needed for running the dev server
make install-tools

# Create an env file with defaults if you haven't already
cp .env.sample .env

# Reset your local storage directories if you haven't already
make reset-local-storage

# Run the dev server
make dev
```
