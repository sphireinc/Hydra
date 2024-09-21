# Use the official Go image as a base
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application files into the container
COPY . .

# Set environment variable for Go path
ENV GO111MODULE=on

# Run tests inside the hydra package
CMD ["go", "test", "-v", "./tests/..."]