# Use golang:1.22 as the base image
FROM golang:1.22

# Set the working directory to /app
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the server application
RUN go build -o server ./server-multi

# Expose port 6660
EXPOSE 6660

# Set the entrypoint to ./server
ENTRYPOINT ["./server"]
