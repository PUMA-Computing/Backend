# Use an official Golang runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required third-party dependencies into the container
RUN go mod tidy

# Build the Go application
RUN go build -o main ./cmd/app

# Expose the port that the application will run on
EXPOSE 8080

# Define the command to run your application
CMD ["./main"]
