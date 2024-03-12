# Use an official Go runtime as a parent image
FROM golang:1.22.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN go mod init modules
RUN go get github.com/go-sql-driver/mysql

# Build the Go application
RUN go build main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ./main