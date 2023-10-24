FROM golang:1.21.3-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

RUN go install github.com/cosmtrek/air@latest

RUN air -v
