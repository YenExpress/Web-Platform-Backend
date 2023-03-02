# Use an official Go runtime as a parent image
FROM golang:1.19-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Download dependencies and build the Go application
RUN go mod download && go build -o main .

# Expose port 8080 for the application to listen on
EXPOSE 8080

# Run the application when the container starts
CMD ["/app/main"]
