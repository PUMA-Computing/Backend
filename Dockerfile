# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app/

# Copy the local package files to the container's workspace
COPY go.mod go.sum /app/

# Download and install any required third-party dependencies into the container
RUN go mod download
COPY . /app/

# Build the Go application
RUN CGO_ENABLED=0 go build \
    -v \
    -trimpath \
    -o main \
    main.go

RUN echo "ID=\"distroless\"" > /etc/os-release

# Stage 2 (Final)
FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/os-release /etc/os-release

COPY --from=builder /app/main /usr/bin/
CMD [ "/usr/bin/main" ]

# Expose the port that the application will run on
EXPOSE 8080