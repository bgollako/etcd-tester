# Use the official Go image as the base image
FROM golang:1.22.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o etcd-tester .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/etcd-tester .

# Use ENTRYPOINT to specify the executable
ENTRYPOINT ["./etcd-tester"]

# Use CMD to provide default arguments (if any)
CMD []
