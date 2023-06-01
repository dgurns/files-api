FROM golang:1.19-alpine AS build

WORKDIR /app

RUN apk add --no-cache make bash

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/bin/files-api /app/cmd/files-api

FROM alpine:3

COPY --from=build /app/bin/files-api /bin/files-api

RUN chmod +x /bin/files-api

EXPOSE 8080

CMD ["/bin/files-api"]