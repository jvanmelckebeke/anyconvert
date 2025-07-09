ARG ALPINE_VERSION=3.19

# Stage 1: Generate swagger docs
FROM golang:1.21.5-alpine${ALPINE_VERSION} as swagger

RUN go install github.com/swaggo/swag/cmd/swag@latest


WORKDIR /app

COPY src/pkg/domain/ .
COPY src/pkg/api/controller.go .
COPY src/main.go .
COPY src/dev.go .

RUN swag init



# Stage 2: Build the Go app
FROM golang:1.21.5-alpine${ALPINE_VERSION} as builder

ENV GOCACHE=/root/.cache/go-build
ENV GIN_MODE=release

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk update && apk add --no-cache \
    git g++ libwebp-dev binutils


# Copy go mod and sum files
COPY src/go.mod src/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN --mount=type=cache,target=/go/src/mod  \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY --from=swagger /app/docs ./docs

COPY src/ .

RUN rm -rf dev.go

ARG FLAGS="prod"

# Build the Go app
RUN --mount=type=cache,target=/go/src/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 \
    GOARCH=$TARGETARCH \
    go build \
    -tags "$FLAGS" \
    -v \
    -ldflags="-s -w" \
    -installsuffix cgo \
    -o main .

FROM alpine:${ALPINE_VERSION} as ffmpeg
RUN apk update

RUN --mount=type=cache,target=/var/cache/apk \
    apk add \
    autoconf \
    automake \
    build-base \
    binutils \
    bzip2 \
    cmake \
    coreutils \
    expat-dev \
    file \
    g++ \
    gcc \
    gperf \
    libtool \
    x264-dev \
    make \
    openssl-dev \
    python3 \
    tar \
    xcb-proto \
    libwebp-dev \
    lame-dev \
    nasm \
    yasm \
    zlib-dev

WORKDIR /app

RUN mkdir -p /app/ffmpeg_sources /app/ffmpeg_build /app/bin /app/lib

WORKDIR /app/ffmpeg_sources

RUN wget -O ffmpeg-snapshot.tar.bz2 https://ffmpeg.org/releases/ffmpeg-snapshot.tar.bz2

RUN tar xjvf ffmpeg-snapshot.tar.bz2

WORKDIR /app/ffmpeg_sources/ffmpeg

ENV PATH="/app/bin:${PATH}"
ENV PKG_CONFIG_PATH="/app/ffmpeg_build/lib/pkgconfig"

RUN ./configure \
    --prefix="/app/ffmpeg_build" \
    --pkg-config-flags="--static" \
    --extra-cflags="-I/app/ffmpeg_build/include" \
    --extra-ldflags="-L/app/ffmpeg_build/lib" \
    --extra-libs="-lpthread -lm" \
    --bindir="/app/bin" \
    --arch=$TARGETARCH \
    \
    --disable-debug \
    --disable-doc \
    --disable-autodetect \
    \
    --enable-gpl \
    --enable-libwebp \
    --enable-libx264 \
    --enable-zlib \
    \
    --disable-ffplay \
    --disable-ffprobe \
    \
    --disable-avdevice \
    --disable-swresample \
    --disable-w32threads \
    --disable-os2threads \
    --disable-network \
    --disable-dwt \
    --disable-lsp \
    --disable-faan \
    --disable-pixelutils \
    --disable-alsa \
    \
    --disable-everything \
    --enable-protocol=file \
    --enable-encoder=libx264,aac,gif,mp3,libmp3lame \
    --enable-decoder=libx264,webp,vp9,vorbis \
    --enable-decoder=gif,png,apng,zlib,aac \
    --enable-filter=pad,scale \
    --enable-muxer=mp4,gif,mov,mp3 \
    --enable-demuxer=gif,mov,image2,matroska,mp3 \
    --enable-parser=png,webp,aac,gif,h264,vorbis,mp3 \
    \
    --enable-libmp3lame




RUN make clean
RUN make -j$(nproc)
RUN make install
RUN hash -r


# fix lib links

CMD ["/bin/sh"]



# Stage 2: Run the Go app
FROM alpine:${ALPINE_VERSION} as runtime


ENV GIN_MODE=release
ENV FFMPEG_VERBOSITY=error

WORKDIR /app/

RUN --mount=type=cache,target=/var/cache/apk \
    apk update && apk add --no-cache \
    libwebp \
    libwebpmux \
    libwebpdemux \
    x264-libs \
    lame-libs \
    \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*

# copy ffmpeg
COPY --from=ffmpeg /app/bin/ffmpeg /app/bin/ffmpeg


# copy app
COPY --from=builder /app/main .


# rebuild links

ENV PATH="/app/bin:${PATH}"

# Expose port 8000 to the outside
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
