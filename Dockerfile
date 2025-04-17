# Use the official Golang image as the build stage
FROM golang:1.22 AS builder

# Set environment variables
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -ldflags "-s -w" -o main .

# Use a Debian base image for the final stage (for LibreOffice)
FROM debian:bookworm-slim

# Install necessary runtime dependencies
RUN apt-get update && \
    apt-get install -y \
    gettext-base \
    libreoffice \
    libreoffice-writer \
    # Fonts for proper rendering
    fonts-liberation \
    fonts-dejavu \
    fonts-freefont-ttf \
    # Clean up to reduce image size
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Create the /storage directory
RUN mkdir -p /storage && chmod -R 777 /storage

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.yaml /app/config.yaml

# Copy the storage directory
COPY storage /app/storage

# Copy any initialization scripts
COPY init-config.sh /app/init-config.sh
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 8002

# Command to run the initialization script
CMD ["/app/init-config.sh"]