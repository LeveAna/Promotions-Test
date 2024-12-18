# Stage 1: Build Stage
FROM golang:1.22 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o promotions-test .

# Stage 2: Run Stage
FROM alpine:3.18

# Install MySQL client for testing connectivity (optional)
RUN apk --no-cache add mysql-client

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=builder /app/promotions-test .

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["./promotions-test"]
