# Use an official Golang runtime as a parent image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Create a directory for binaries
RUN mkdir -p /app/binaries
# Build the main application
# RUN go build -o todoapp && ls -l /app
# If there are multiple binaries to build, repeat the RUN command with different paths
RUN go build -o binaries/todoapp ./cmd && ls -l /app
RUN go build -o binaries/workerapp ./consumer && ls -l /app


# Build the Faktory worker
# RUN go build -o faktory_worker ./cmd/worker && ls -l /app

# Expose the port that your app will run on
EXPOSE 8080

# Set the working directory for runtime
WORKDIR /app/binaries

# Command to run the main application by default
CMD ["./todoapp"]
