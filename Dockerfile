# Use the official Golang image as a build stage
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod .
COPY go.sum .

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go mod download
RUN go build -o main .

# Use a minimal Alpine Linux image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the application
CMD ["./main"]
