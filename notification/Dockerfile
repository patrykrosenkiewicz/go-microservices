# Step 1: Use the official Golang image to build the Go application
FROM golang:1.23-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application source code into the container
COPY . .

# Step 2: Build the Go application inside the container
RUN go build -o main .

# Step 3: Create the final image using a smaller image
FROM alpine:latest

# Install CA certificates to allow HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory inside the final image
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Expose the port on which the app will run
EXPOSE 8080

# Command to run the Go application when the container starts
CMD ["./main"]