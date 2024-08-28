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

# Build the main application
RUN go build -o todoapp && ls -l /app

# Build the Faktory worker
# RUN go build -o faktory_worker ./cmd/worker && ls -l /app

# Expose the port that your app will run on
EXPOSE 8080

# Command to run the main application by default
CMD ["./todoapp"]
