# syntax=docker/dockerfile:1

# ========== BUILD ==========
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.1

COPY go.mod ./

RUN go mod download

COPY . .

CMD ["air"]
