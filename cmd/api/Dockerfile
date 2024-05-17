FROM golang:1.20.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the containder's workspace
COPY . .

# Build the Go app
RUN go build -o main ./cmd/api

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]