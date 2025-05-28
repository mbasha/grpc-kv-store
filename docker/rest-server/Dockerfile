# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the REST API server executable
# CGO_ENABLED=0 is important for static binaries
# -o specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kvstore-rest-server cmd/kvstore-rest-server/main.go

# Stage 2: Create the final lean image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled executable from the builder stage
COPY --from=builder /app/kvstore-rest-server .

# Expose the HTTP port
EXPOSE 8080

# Command to run the executable
CMD ["./kvstore-rest-server"]