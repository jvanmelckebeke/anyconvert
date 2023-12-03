ARG BASE_TAG_ARM64=3.10.0-buster
ARG BASE_TAG_AMD64=3.10.0-alpine3.15

ARG RUNTIME_TAG_ARM64=3.10.0-slim-buster
ARG RUNTIME_TAG_AMD64=3.10.0-alpine3.15


FROM python:${BASE_TAG_ARM64} as base-arm64

ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    libssl-dev \
    libffi-dev \
    python3-dev \
    python3-venv \
    python3-pip


FROM python:${BASE_TAG_AMD64} as base-amd64

ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apk update && \
    apk add \
    build-base \
    gcc \
    g++ \
    gfortran \
    openblas-dev

FROM base-${TARGETARCH} as build

ENV VIRTUAL_ENV=/venv
ENV venv_dir=/venv

RUN python -m venv $VIRTUAL_ENV

ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# for compatibility with buildx caches
ENV XDG_CACHE_HOME=/root/.cache

# install wheel for faster builds
RUN --mount=type=cache,target=/root/.cache/,sharing=locked \
    pip install --upgrade pip setuptools wheel

RUN python -m venv $VIRTUAL_ENV

COPY requirements.txt requirements.txt

ENV PIP_OPTIONS="--timeout=100"

# install dependencies and clean pycache files
RUN --mount=type=cache,target=/root/.cache/,sharing=locked \
    pip install $PIP_OPTIONS -r requirements.txt && \
    find /venv -name "*.pyc" -delete && \
    find /venv -name "__pycache__" -exec rm -rf {} +

FROM debian:bookworm as imagemagick

ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apt-get update

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apt-get install -y --no-install-recommends \
    wget build-essential cmake openssl git ca-certificates

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apt-get install -y --no-install-recommends \
    git curl make cmake automake libtool yasm g++ \
    pkg-config perl libde265-dev libx265-dev libltdl-dev \
    libopenjp2-7-dev liblcms2-dev libbrotli-dev libzip-dev \
    libbz2-dev liblqr-1-0-dev libzstd-dev libgif-dev libjpeg-dev \
    libopenexr-dev libpng-dev libwebp-dev librsvg2-dev libwmf-dev \
    libxml2-dev libtiff-dev libraw-dev ghostscript gsfonts ffmpeg \
    libpango1.0-dev libdjvulibre-dev libfftw3-dev libgs-dev libgraphviz-dev


RUN mkdir -p  /magick /magick-config /magick-build /magick-work


RUN wget https://dist.1-2.dev/imei.sh -qO /imei.sh && \
    chmod +x /imei.sh


RUN bash ./imei.sh \
    --no-backports \
    --skip-jxl \
    --skip-libheif \
    --build-dir /magick

# dummy entrypoint
ENTRYPOINT ["/bin/bash"]

FROM python:${RUNTIME_TAG_ARM64} as dependencies-arm64

ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=cache,target=/var/cache/apt/,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    ffmpeg \
    && \
    rm -rf /var/lib/apt/lists/* \
    && \
    apt-get clean

FROM python:${RUNTIME_TAG_AMD64} as dependencies-amd64

ENV DEBIAN_FRONTEND=noninteractive

RUN apk update && \
  apk add --no-cache\
    ffmpeg \
#    imagemagick \
    && rm -rf /var/cache/apk/*

FROM dependencies-${TARGETARCH} as dependencies

COPY --from=imagemagick /magick/bin /usr/local/bin
COPY --from=imagemagick /magick/lib /usr/local/lib
COPY --from=imagemagick /magick/include /usr/local/include
#COPY --from=imagemagick /magick/share /usr/local/share

FROM dependencies as runtime

ENV PYTHONUNBUFFERED=1
ENV GRADIO_PORT=8080

ENV VIRTUAL_ENV=/venv
ENV VENV_DIR=/venv


RUN mkdir -p /tmp/gradio

COPY --from=build /venv $VIRTUAL_ENV

ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV PYTHONPATH="$VIRTUAL_ENV/lib/python3.10/site-packages:$PYTHONPATH"

# otherwise matplotlib will always try to build the font cache
RUN python -c "import matplotlib.pyplot"

WORKDIR /app

COPY main.py .
COPY src/ src/

CMD ["python", "main.py"]



