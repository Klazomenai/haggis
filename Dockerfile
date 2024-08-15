# Use an official Go runtime as a parent image
FROM golang:1.23.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the host machine into the container
COPY . .

# Build the Go application
RUN go build -o haggis .

# Set the default command to run the haggis application
ENTRYPOINT ["./haggis"]
