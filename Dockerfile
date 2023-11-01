FROM python:3.10.0-alpine3.15

ENV PYTHONUNBUFFERED=1
ENV GRADIO_PORT=8080

RUN apk update
RUN apk add --no-cache ffmpeg

RUN mkdir -p /tmp/gradio

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . .

CMD ["python", "main.py"]
