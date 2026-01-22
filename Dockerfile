# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Install required runtime packages (tzdata for timezone support)
RUN apk add --no-cache tzdata

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy migrations
COPY migrations ./migrations

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
