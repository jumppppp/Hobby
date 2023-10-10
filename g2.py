import argparse
import time
# 创建命令行参数解析器
parser = argparse.ArgumentParser(description='Copy the contents of an input file to an output file.')

# 添加-s参数，用于指定输入文件
parser.add_argument('-s', '--source', required=True, help='Input file path')

# 添加-o参数，用于指定输出文件
parser.add_argument('-o', '--output', required=True, help='Output file path')

# 解析命令行参数
args = parser.parse_args()

# 打开输入文件并读取内容
try:
    with open(args.source, 'r') as source_file:
        file_content = source_file.read()
        time.sleep(5)
    # 打开输出文件并将内容写入

    with open(args.output, 'w') as output_file:
        output_file.write(file_content)

    print(f'File copied from {args.source} to {args.output}')
except FileNotFoundError:
    print('Input file not found.')
except Exception as e:
    print(f'An error occurred: {str(e)}')