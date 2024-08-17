# Start from a small base image
FROM golang:1.23.0-alpine3.20 as builder

# Install build dependencies
RUN apk add --no-cache build-base

# Set the working directory
WORKDIR /app

# Copy only the necessary Go files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 \
    go build -tags "modern,migrate" -o app ./cmd/app/

# Create a minimal runtime image
FROM alpine:3.20.2

# Install runtime dependencies
RUN apk add --no-cache tzdata

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

COPY db/migrations ./migrations

ENV STORAGE_MIGRATIONS_PATH=/app/migrations

ENTRYPOINT [ "./app" ]
