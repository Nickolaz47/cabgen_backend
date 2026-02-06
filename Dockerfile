# Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Test Dependecies
RUN apk add --no-cache gcc musl-dev

# Source code
COPY internal/ ./internal
COPY cmd/ ./cmd
COPY . .

# Tests
RUN go test -v ./...

# Compilation
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/server

# Runtime
FROM alpine:latest

WORKDIR /app

# Getting the binary from builder
COPY --from=builder /app/api .

COPY --from=builder /app/internal/translation/active ./internal/translation/active
COPY --from=builder /app/internal/data ./internal/data

EXPOSE 8080

RUN mkdir ./logs

ENTRYPOINT ["./api"]