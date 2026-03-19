# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy rest
COPY . .

# Build binary
RUN go build -o app ./cmd/main/main.go


# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY config.yaml .

EXPOSE 8080

CMD ["./app"]
