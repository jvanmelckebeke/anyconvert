from PIL import Image

from src.tools import convert_input_to_output_path


def webp_to_jpg(input):
    if input is None:
        return None

    output_filename = convert_input_to_output_path(input, "jpg")

    image = Image.open(input).convert("RGB")
    image.save(output_filename, format="jpeg")

    print(f"saved image to {output_filename}")

    return output_filename
