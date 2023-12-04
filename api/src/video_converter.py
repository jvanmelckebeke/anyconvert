import subprocess
import time

from src.tools import convert_input_to_output_path, filesize, bytes2human, convert_input_to_output_directories


def ffmpeg_process(input_file_path, output_file_path, subprocess_args, progress=None):
    with subprocess.Popen(subprocess_args) as process:
        raw_input_size = filesize(input_file_path)
        input_size = bytes2human(raw_input_size)

        while process.poll() is None:
            print("waiting for ffmpeg to finish...")
            try:
                raw_output_size = filesize(output_file_path)
            except FileNotFoundError:
                raw_output_size = 0

            output_size = bytes2human(raw_output_size)

            print(f"input size: {input_size}, output size: {output_size}")

            if progress is not None:
                progress((raw_output_size, raw_input_size), unit="bytes")

            time.sleep(1)

        if process.returncode != 0:
            raise Exception("ffmpeg failed to convert the file")
    return output_file_path


def media_to_mp4(input, progress=None):
    if input is None:
        return None

    if input.endswith(".webp"):
        return animated_webp_to_mp4(input, progress)
    elif input.endswith(".webm"):
        return webm_to_mp4(input, progress)
    elif input.endswith(".gif"):
        return webm_to_mp4(input, progress)
    else:
        raise Exception("unsupported file format")


def webm_to_mp4(input, progress):
    input_file_path = input
    output_file_path = convert_input_to_output_path(input_file_path, "mp4")

    # run ffmpeg to convert the file to mp4
    # the command to be run is
    subprocess_args = [
        "ffmpeg",
        "-y",  # overwrite output file if it exists
        "-v", "error",
        "-i", input_file_path,
        "-map", "V:0?",
        "-map", "0:a?",
        "-c:v", "libx264",
        "-movflags", "+faststart",
        "-preset", "veryslow",
        "-pix_fmt", "yuv420p",
        "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
        output_file_path
    ]

    return ffmpeg_process(input_file_path, output_file_path, subprocess_args, progress)


def animated_webp_to_mp4(input, progress):
    input_file_path = input
    output_file_path = convert_input_to_output_path(input_file_path, "mp4")

    # first convert the webp to a png sequence
    png_sequence_path = convert_input_to_output_directories(input_file_path)
    subprocess_args = [
        "magick",
        input_file_path,
        "-coalesce",  # convert the webp to a png sequence
        f"{png_sequence_path}/frames.png"
    ]
    with subprocess.Popen(subprocess_args) as magick_process:
        while magick_process.poll() is None:
            print("waiting for magick to finish...")
            time.sleep(1)

        if magick_process.returncode != 0:
            raise Exception("magick failed to convert the file")

    # then convert the png sequence to mp4
    subprocess_args = [
        "ffmpeg",
        "-y",  # overwrite output file if it exists
        "-v", "error",
        "-i", f"{png_sequence_path}/frames-%d.png",
        "-c:v", "libx264",
        "-movflags", "+faststart",
        "-preset", "veryslow",
        "-pix_fmt", "yuv420p",
        "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
        output_file_path
    ]

    return ffmpeg_process(input_file_path, output_file_path, subprocess_args, progress)
