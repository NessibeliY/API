FROM golang:1.20.1-alpine


# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download all the dependencies
RUN go mod download

# Build the Go app
RUN go build -o api ./cmd/api

# Expose port 8081 to the outside world
EXPOSE 4000

# Command to run the executable with wait-for-it
CMD ["./api"]
