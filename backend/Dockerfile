# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Stage 2: Serve
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests (Groq API)
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary
COPY --from=builder /app/main .

# Copy the MaxMind DB and .env file (if running locally)
COPY GeoLite2-City.mmdb .
COPY .env* ./

EXPOSE 8080

CMD ["./main"]
