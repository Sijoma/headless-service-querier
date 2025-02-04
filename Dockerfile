# Use the official Golang image as the base
FROM golang:1.23.5 AS builder
ARG TARGETOS
ARG TARGETARCH
# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.mod
RUN go mod download

# Copy the rest of the application source code
COPY main.go main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o headless

# Use a lightweight image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/headless .

# Expose port (if needed)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/headless"]
