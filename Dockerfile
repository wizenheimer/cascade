# Start from the latest golang base image
FROM golang:1.22.5-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY src ./src

# Set the working directory to where go.mod is located
WORKDIR /app/src

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Development stage with hot reload
FROM golang:1.22.5-alpine AS development

# Install git and any other necessary tools
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the entire src directory and .air.toml file to the current working directory
COPY . .

# Copy the entire src directory
COPY src ./src

# Copy .air.toml file to the root of the app directory
COPY .air.toml /app/.air.toml

# Set the working directory to where go.mod is located
WORKDIR /app/src

# Download all dependencies
RUN go mod download

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the development server with hot reloading
CMD ["air", "-c", "/app/.air.toml"]

# Start a new stage from scratch for production
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/src/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
