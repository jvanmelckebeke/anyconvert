# Stage 1: Build the Go app
FROM golang:1.21.5-alpine3.18 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk update && apk add --no-cache \
    git libwebp-dev gcc g++ make

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY src/ .

# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 GOOS=linux go build -a -o main .

# Stage 2: Run the Go app
FROM alpine:latest

RUN apk --no-cache add libwebp-dev

WORKDIR /root/


# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000 to the outside
EXPOSE 8000

# Command to run the executable
CMD ["./main"]