FROM python:3.10.0-alpine3.15

ENV PYTHONUNBUFFERED=0

RUN apk update
RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . .

CMD ["python", "main.py"]
