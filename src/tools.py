from hurry.filesize import size, alternative
import os


def filesize(filename):
    return os.stat(filename).st_size


def bytes2human(raw):
    return size(raw, system=alternative)


def convert_input_to_output_path(input_fname, output_ext):
    basename_with_ext = os.path.basename(input_fname)
    basename, _ = os.path.splitext(basename_with_ext)

    return f"/tmp/{basename}.{output_ext}"


def convert_input_to_output_directories(input_fname):
    basename_with_ext = os.path.basename(input_fname)
    basename, _ = os.path.splitext(basename_with_ext)

    if not os.path.exists(f"/tmp/{basename}/"):
        os.mkdir(f"/tmp/{basename}/")
    else:
        # remove all files in the directory
        for f in os.listdir(f"/tmp/{basename}/"):
            os.remove(os.path.join(f"/tmp/{basename}/", f))

    return f"/tmp/{basename}/"
