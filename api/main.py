#!/bin/env python3
import os
import asyncio
import uuid
from datetime import datetime

from fastapi import FastAPI, UploadFile, File, HTTPException, BackgroundTasks
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

UPLOADS_DIR = "/tmp"

uploads = {}


def get_filepath(filesource: str):
    return f"{UPLOADS_DIR}/{filesource}"


def save_file_and_create_upload_id(file: UploadFile):
    uid = str(uuid.uuid4())[0:8]
    uploads[uid] = {
        "upload_id": uid,
        "filename": file.filename,
        "filesource": f"{uid}_{file.filename}",
        "created_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "status": "pending"
    }

    filepath = get_filepath(uploads[uid]["filesource"])
    with open(filepath, "wb") as buffer:
        buffer.write(file.file.read())

    print(f"File saved to {filepath}")
    return uid


async def delete_upload(upload_id: str, delay: int):
    await asyncio.sleep(delay)
    try:
        os.remove(uploads[upload_id])
    except FileNotFoundError:
        pass
    finally:
        del uploads[upload_id]


async def process_upload(upload_id: str, converter_func):
    try:
        uploads[upload_id]["status"] = "processing"
        filepath = get_filepath(uploads[upload_id]["filesource"])
        output_path = converter_func(filepath)

        # remove the UPLOADS_DIR from the output path
        if output_path.startswith(UPLOADS_DIR):
            output_path = output_path[len(UPLOADS_DIR) + 1:]

        uploads[upload_id]["output_path"] = output_path
        uploads[upload_id]["result_url"] = f"/uploads/{upload_id}"
        uploads[upload_id]["status"] = "done"
    except Exception as e:
        uploads[upload_id]["status"] = "error"
        uploads[upload_id]["error"] = str(e)
        raise e
    finally:
        # delete the source file
        try:
            filepath = get_filepath(uploads[upload_id]["filesource"])
            os.remove(filepath)
        except FileNotFoundError:
            pass
        # after 5 minutes, delete the file
        asyncio.create_task(delete_upload(upload_id, 300))


@app.get("/status")
async def get_uploads():
    return uploads


@app.get("/uploads/{upload_id}")
async def get_file(upload_id: str):
    if upload_id not in uploads:
        raise HTTPException(status_code=404, detail="Upload not found")

    filepath = get_filepath(uploads[upload_id]["output_path"])
    return FileResponse(filepath)


@app.post("/image")
async def upload_image(file: UploadFile = File(...), background_tasks: BackgroundTasks = None):
    upload_id = save_file_and_create_upload_id(file)

    background_tasks.add_task(process_upload, upload_id, webp_to_jpg)

    return {"status": f"/status/{upload_id}"}


@app.post("/video")
async def upload_video(file: UploadFile = File(...), background_tasks: BackgroundTasks = None):
    upload_id = save_file_and_create_upload_id(file)

    background_tasks.add_task(process_upload, upload_id, media_to_mp4)

    return {"status": f"/status/{upload_id}"}


@app.get("/status/{upload_id}")
async def get_upload_status(upload_id: str):
    if upload_id not in uploads:
        raise HTTPException(status_code=404, detail="Upload not found")
    return uploads[upload_id]


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=PORT)
