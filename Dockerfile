FROM python:3.10.0-alpine3.15 as build

ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=cache,target=/var/cache/apt \
    apk update && \
    apk add \
    build-base \
    gcc \
    g++ \
    gfortran \
    openblas-dev

ENV VIRTUAL_ENV=/venv
ENV venv_dir=/venv

RUN python -m venv $VIRTUAL_ENV

ENV PATH="$VIRTUAL_ENV/bin:$PATH"
# for compatibility with buildx caches
ENV XDG_CACHE_HOME=/root/.cache

COPY requirements.txt requirements.txt

ENV PIP_OPTIONS="--timeout=100"

RUN --mount=type=cache,target=/root/.cache/ pip install $PIP_OPTIONS -r requirements.txt && \
    find /venv -name "*.pyc" -delete && \
    find /venv -name "__pycache__" -exec rm -rf {} + && \
#    find /venv -name "*.dist-info" -exec rm -rf {} + && \
    rm -rf /root/.cache/pip
FROM python:3.10.0-alpine3.15 as runtime

ENV PYTHONUNBUFFERED=1
ENV GRADIO_PORT=8080

ENV VIRTUAL_ENV=/venv
ENV VENV_DIR=/venv

RUN apk update && \
  apk add --no-cache ffmpeg imagemagick && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /tmp/gradio

WORKDIR /app

COPY --from=build /venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV PYTHONPATH="$VIRTUAL_ENV/lib/python3.10/site-packages:$PYTHONPATH"


COPY main.py .
COPY src/ src/

CMD ["python", "main.py"]
