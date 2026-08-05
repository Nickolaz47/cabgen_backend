# Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source code
COPY internal/ ./internal
COPY cmd/ ./cmd
COPY . .

# Tests
RUN adduser -D testuser && chown -R testuser /app
USER testuser
RUN CGO_ENABLED=0 go test -v ./...

# Compilation
USER root
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker-email ./cmd/worker-email

# Runtime
FROM scratch

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/worker-email .

COPY --from=builder /app/internal/translation/active ./internal/translation/active
COPY --from=builder /app/jsons ./jsons

EXPOSE 8080

CMD ["./api"]