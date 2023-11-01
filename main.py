import os
import time

import gradio as gr
import subprocess

import io
from PIL import Image
from hurry.filesize import size, alternative


def filesize(filename):
    return os.stat(filename).st_size


def bytes2human(raw):
    return size(raw, system=alternative)


def webp_to_jpg(input):
    if input is None:
        return None

    image = Image.open(input).convert("RGB")
    output = io.BytesIO()
    image.save(output, format="jpeg")
    return Image.open(output)


def webm_to_mp4(input, progress=gr.Progress()):
    input_file_path = input

    if input_file_path is None:
        return None

    basename_with_ext = os.path.basename(input_file_path)
    basename, _ = os.path.splitext(basename_with_ext)

    output_file_path = f"/tmp/{basename}.mp4"

    # run ffmpeg to convert the file to mp4
    # the command to be run is
    # ffmpeg -v error -map V:0? -map 0:a? -c:v libx264 -movflags +faststart -preset veryslow -vf pad=ceil(iw/2)*2:ceil(ih/2)*2 -i /tmp/input.webm
    ffmpeg_process = subprocess.Popen([
        "ffmpeg",
        "-y",  # overwrite output file if it exists
        "-v", "error",
        "-i", input_file_path,
        "-map", "V:0?",
        "-map", "0:a?",
        "-c:v", "libx264",
        "-movflags", "+faststart",
        "-preset", "veryslow",
        "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
        output_file_path
    ])

    raw_input_size = filesize(input_file_path)
    input_size = bytes2human(raw_input_size)

    while ffmpeg_process.poll() is None:
        print("waiting for ffmpeg to finish...")
        try:
            raw_output_size = filesize(output_file_path)
        except FileNotFoundError:
            raw_output_size = 0

        output_size = bytes2human(raw_output_size)

        print(f"input size: {input_size}, output size: {output_size}")

        progress_ratio = raw_output_size / raw_input_size

        progress((raw_output_size, raw_input_size), unit="bytes")

        time.sleep(1)

    if ffmpeg_process.returncode != 0:
        raise Exception("ffmpeg failed to convert the file")

    return "/tmp/output.mp4"


with gr.Blocks() as iface:
    with gr.Tab("webp to jpg"):
        with gr.Row():
            with gr.Column(scale=1):
                in_image = gr.File(label="Input", file_types=[".webp"])
            with gr.Column(scale=1):
                out_image = gr.Image(label="Output")
        with gr.Row():
            img_submit_btn = gr.Button("Submit", variant="primary")

        img_submit_btn.click(fn=webp_to_jpg, inputs=in_image, outputs=out_image)
    with gr.Tab("webm to mp4"):
        with gr.Row():
            with gr.Column(scale=1):
                in_video = gr.File(label="Input", file_types=[".webm"])
            with gr.Column(scale=1):
                out_video = gr.Video(label="Output")

        with gr.Row():
            submit_btn = gr.Button("Submit", variant="primary")

        submit_btn.click(fn=webm_to_mp4, inputs=in_video, outputs=out_video)

if __name__ == '__main__':
    iface.launch()
