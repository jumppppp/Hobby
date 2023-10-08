import os
import argparse
import math

def get_file_size(file_path):
    try:
        size = os.path.getsize(file_path)
        return convert_size(size)
    except FileNotFoundError:
        return "File not found"

def convert_size(size_bytes):
    if size_bytes == 0:
        return "0B"
    size_name = ("B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB")
    i = int(math.floor(math.log(size_bytes, 1024)))
    p = math.pow(1024, i)
    s = round(size_bytes / p, 2)
    return f"{s} {size_name[i]}"

def main():
    parser = argparse.ArgumentParser(description='Get the size of a file.')
    parser.add_argument('-f', '--file', help='File to be measured', required=True)
    args = parser.parse_args()

    file_path = args.file
    file_size = get_file_size(file_path)
    
    with open('file_size_output.txt', 'w') as f:
        f.write(f"The size of {file_path} is: {file_size}")

if __name__ == '__main__':
    main()