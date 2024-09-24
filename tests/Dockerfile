# Use the official Go image as a base
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from the project root
COPY ../go.mod ./
COPY ../go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project from the context root, including the hydra package and tests
COPY ../ ./

# Set environment variable for Go path
ENV GO111MODULE=on

# Run tests inside the hydra package
CMD ["go", "test", "-v", "./tests/..."]