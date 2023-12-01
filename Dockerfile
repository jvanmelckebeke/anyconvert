FROM python:3.10.0-buster as build

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    ninja-build \
    g++ \
    gcc \
    libc-dev

ENV VIRTUAL_ENV=/venv
ENV venv_dir=/venv

RUN python -m venv $VIRTUAL_ENV

ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# for compatibility with buildx caches
ENV XDG_CACHE_HOME=/root/.cache

RUN --mount=type=cache,target=/root/.cache/ python -m pip install --upgrade pip wheel

COPY requirements.txt requirements.txt

RUN --mount=type=cache,target=/root/.cache/ pip install -r requirements.txt

FROM alpine:3.15.0 as runtime

ENV PYTHONUNBUFFERED=1
ENV GRADIO_PORT=8080

ENV VIRTUAL_ENV=/venv
ENV venv_dir=/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

RUN apk update && \
  apk add --no-cache ffmpeg imagemagick

RUN mkdir -p /tmp/gradio

WORKDIR /app

COPY --from=build /venv /venv

COPY main.py .
COPY src/ src/

CMD ["python", "main.py"]
