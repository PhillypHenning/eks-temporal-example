# Start from the official Golang image to build the application
FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create and set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files if you have dependencies
COPY go.mod go.sum ./

# Download necessary Go modules
RUN go mod download

# Copy the entire app directory contents
COPY . .

# Build the executable binary
RUN go build -o worker .

# Create a small image from Alpine Linux
FROM alpine:3.18

# Copy the binary from the builder stage
COPY --from=builder /app/worker /worker

# Expose any ports the app requires
EXPOSE 8080

# Command to run the executable
CMD ["/worker"]