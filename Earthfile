VERSION 0.7
FROM golang:1.21-alpine
WORKDIR /duet

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

lint:
    FROM golangci/golangci-lint:latest
    COPY . .
    RUN golangci-lint run

build:
    FROM +deps
    COPY . .
    RUN CGO_ENABLED=0 go build -o duet cmd/duet/main.go
    SAVE ARTIFACT duet AS LOCAL dist/duet

test:
    FROM +deps
    COPY . .
    RUN go test -race -coverprofile=coverage.out ./...
    SAVE ARTIFACT coverage.out AS LOCAL coverage.out

docker:
    FROM alpine:latest
    COPY +build/duet /usr/local/bin/duet
    ENTRYPOINT ["/usr/local/bin/duet"]
    SAVE IMAGE duet:latest

all:
    BUILD +lint
    BUILD +test
    BUILD +build
    BUILD +docker
