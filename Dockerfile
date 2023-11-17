FROM python:3.10.0-alpine3.15 as base

ENV PYTHONUNBUFFERED=1
ENV GRADIO_PORT=8080

RUN apk update && apk add --no-cache ffmpeg imagemagick

RUN mkdir -p /tmp/gradio

WORKDIR /app

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt --no-cache-dir --timeout 1000 && \
  rm -rf /root/.cache/pip

COPY main.py .
COPY src/ src/

CMD ["python", "main.py"]
