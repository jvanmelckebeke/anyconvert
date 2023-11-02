#!/bin/env python3
import os

import gradio as gr

from src.image_converter import webp_to_jpg
from src.video_converter import media_to_mp4

with gr.Blocks(title="Webany converter") as iface:
    with gr.Tab("webp to jpg"):
        with gr.Row():
            with gr.Column(scale=1):
                in_image = gr.File(label="Input", file_types=[".webp"])
            with gr.Column(scale=1):
                out_image = gr.Image(label="Output")
        with gr.Row():
            img_submit_btn = gr.Button("Submit", variant="primary")

        img_submit_btn.click(
            fn=webp_to_jpg, inputs=in_image, outputs=out_image)
    with gr.Tab("to mp4"):
        with gr.Row():
            with gr.Column(scale=1):
                in_video = gr.File(label="Input", file_types=[
                                   ".webm", ".webp", ".gif"])
            with gr.Column(scale=1):
                out_video = gr.Video(label="Output")

        with gr.Row():
            submit_btn = gr.Button("Submit", variant="primary")

        submit_btn.click(fn=media_to_mp4, inputs=[in_video], outputs=out_video)

if __name__ == '__main__':
    port = int(os.environ.get("GRADIO_PORT", 8088))

    iface.launch(server_port=port, server_name="0.0.0.0")
