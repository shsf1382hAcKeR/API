# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the local code to the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["go", "run", "./cmd/server/main.go"]
