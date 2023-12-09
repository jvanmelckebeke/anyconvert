# Stage 1: Generate swagger docs
FROM golang:1.21.5-alpine3.18 as swagger

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY src/domain/ .
COPY src/api/ .
COPY src/main.go .


RUN swag init

# Stage 2: Build the Go app
FROM golang:1.21.5-alpine3.19 as builder

ENV GOCACHE=/root/.cache/go-build
ENV GIN_MODE=release

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk update && apk add --no-cache \
    git g++ libwebp-dev


# Copy go mod and sum files
COPY src/go.mod src/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/pkg/mod  \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY --from=swagger /app/docs ./docs

COPY src/ .

# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=$TARGETARCH \
    go build \
    -a \
    -ldflags="-s -w" \
    -linkshared \
    -installsuffix cgo \
    -o main .

# Stage 2: Run the Go app
FROM alpine:3.17 as runtime


ENV GIN_MODE=release

WORKDIR /app/

RUN --mount=type=cache,target=/var/cache/apk \
    apk update && apk add --no-cache \
    ffmpeg \
    libwebp-dev \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000 to the outside
EXPOSE 8000

# Command to run the executable
CMD ["./main"]