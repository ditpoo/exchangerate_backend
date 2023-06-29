# Use an official Golang runtime as the base image
FROM golang:1.20-alpine

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Download all the dependencies
# RUN go mod download
RUN source ./setup.sh

# Build the Go application
RUN go build -o app

# Expose the desired port (change if necessary)
EXPOSE 8080

# Set the command to run the executable
CMD ["./app"]
