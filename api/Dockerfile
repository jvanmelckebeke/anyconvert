
FROM python:3.10.0-alpine3.15 as builder

RUN apk update && apk add \
    build-base g++ gcc gfortran openblas-dev


ENV VIRTUAL_ENV=/venv
ENV VENV_DIR=/venv


ENV PATH="$VIRTUAL_ENV/bin:$PATH"

ENV XDG_CACHE_HOME=/root/.cache/

RUN --mount=type=cache,target=/root/.cache/pip,sharing=locked \
    pip install --upgrade pip setuptools wheel

RUN python -m venv $VIRTUAL_ENV

COPY requirements.txt .

RUN --mount=type=cache,target=/root/.cache/pip,sharing=locked \
    pip install -r requirements.txt

ENV PIP_OPTIONS="--timeout=100"

RUN --mount=type=cache,target=/root/.cache/pip,sharing=locked \
    pip install -r requirements.txt && \
    find /venv -name "*.py[cod]" -delete && \
    find /venv -name "__pycache__" -exec rm -rf {} +


# rm all packages that are not needed
WORKDIR /venv/lib/python3.10/site-packages

RUN find . -name '__pycache__' -type d -exec rm -rf {} +

RUN find . -name '*.pyc' -delete && \
    find . -name '*.pyo' -delete && \
    find . -name '*~' -delete

RUN find . -name '*.dist-info' -type d -exec rm -rf {} +

RUN find . -name 'ensurepip' -type d -exec rm -rf {} +
RUN find . -name 'pip' -type d -exec rm -rf {} +

FROM alpine:3.15 as python-builder

RUN --mount=type=cache,target=/root/.cache/apk,sharing=locked \
    apk update && apk add \
    build-base g++ gcc gfortran openblas-dev wget tar xz make cmake curl libbz2 zlib-dev openssl-dev ncurses-dev readline-dev sqlite-dev tk-dev libffi-dev xz-dev


RUN mkdir -p '/opt/python'

WORKDIR /build

RUN wget https://www.python.org/ftp/python/3.10.0/Python-3.10.0.tgz

RUN tar -xzf Python-3.10.0.tgz

WORKDIR /build/Python-3.10.0

RUN ./configure \
    --with-ensurepip=no \
    --without-doc-strings \
    --without-mimalloc \
    --without-pymalloc \
    --without-readline \
    --without-ensurepip \
    --without-ffi \
    --with-lto \
    --with-system-ffi \
    --disable-test-modules \
    --prefix=/opt/python
RUN make -j "$(nproc)" LDFLAGS="-Wl,--strip-all"
RUN make altinstall



# remove everything we do not need

WORKDIR /opt/python


RUN find . -depth \( \( -type d -a \( -name test -o -name tests -o -name idle_test \) \) -o \( -type f -a \( -name '*.pyc' -o -name '*.pyo' -o -name '*.a' \) \) \) -exec rm -rf '{}' +

RUN find . -name '__pycache__' -type d -exec rm -rf {} +

RUN find . -name '*.pyc' -delete && \
    find . -name '*.pyo' -delete && \
    find . -name '*~' -delete

RUN find . -name 'ensurepip' -type d -exec rm -rf {} +

RUN find . -name 'config-3.10-x86_64-linux-gnu' -type d -exec rm -rf {} +


# link python
RUN ln -s "/opt/python/bin/python3.10" "/opt/python/bin/python"

FROM alpine:3.15 as dependencies

RUN --mount=type=cache,target=/root/.cache/apk,sharing=locked \
    apk update && apk add --no-cache \
    ffmpeg \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /var/cache/apk/*


COPY --from=python-builder /opt/python /opt/python

ENV PYTHONPATH="/opt/python/lib/python3.10/site-packages"


ENV PATH="/opt/python/bin:$PATH"


COPY --from=builder /venv/lib/python3.10/site-packages /opt/python/lib/python3.10/site-packages

FROM dependencies as production

ENV PYTHONUNBUFFERED=1
ENV PORT="8000"

WORKDIR /app

COPY main.py .
COPY src/ src/

CMD ["sh", "-c", "python -m uvicorn main:app --host '0.0.0.0' --port $PORT"]