#!/bin/env python3
import os

from fastapi import FastAPI, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse

from src.image_converter import webp_to_jpg
from src.video_converter import media_to_mp4

PORT = os.getenv("PORT", 8000)

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/uploads/{filename}")
async def get_file(filename: str):
    return FileResponse(f"/tmp/{filename}")


@app.post("/image")
async def upload_image(file: UploadFile = File(...)):
    identifier = file.filename.split(".")[0]
    with open(f"/tmp/{file.filename}", "wb") as buffer:
        buffer.write(file.file.read())

    new_file = webp_to_jpg(f"/tmp/{file.filename}")

    new_file = new_file.replace("/tmp", "")

    return {"file": f"/uploads{new_file}"}


@app.post("/video", response_class=FileResponse)
async def upload_video(file: UploadFile = File(...)):
    with open(f"/tmp/{file.filename}", "wb") as buffer:
        buffer.write(file.file.read())

    new_file = media_to_mp4(f"/tmp/{file.filename}")

    print(new_file)

    return FileResponse(new_file, media_type="video/mp4")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=PORT)
