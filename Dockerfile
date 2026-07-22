# Build
FROM golang:1.25-alpine AS builder

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
RUN adduser -D testuser && chown -R testuser /app
USER testuser
RUN go test -v ./...

# Compilation
USER root
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/server
RUN CGO_ENABLED=1 GOOS=linux go build -o worker-email ./cmd/worker-email

# Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/worker-email .

COPY --from=builder /app/internal/translation/active ./internal/translation/active
COPY --from=builder /app/jsons ./jsons

EXPOSE 8080

RUN mkdir ./logs

CMD ["./api"]