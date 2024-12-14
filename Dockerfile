# syntax=docker/dockerfile:1

# ========== BUILD ==========
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/logparser && ls -la /bin/app

# ========== PROD ==========
FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/app /bin/app

RUN ls -la /bin/app && chmod +x /bin/app

ENTRYPOINT ["/bin/app"]