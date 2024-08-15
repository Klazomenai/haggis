# Use the official Golang image as the base image
FROM golang:1.23.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the entire project source to the working directory
COPY . .

# Build the Go application
RUN go build -o haggis main.go
