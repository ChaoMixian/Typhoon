#!/bin/bash

# 设置项目根目录
PROJECT_DIR=$(pwd)

# 设置目标构建目录（可以是测试目录）
BUILD_DIR="${PROJECT_DIR}/build"
TEST_DIR="${PROJECT_DIR}/tests"

# 设置二进制文件名
BINARY_NAME="Typhoon"  # 根据你的项目名称修改
BINARY_PATH="${BUILD_DIR}/${BINARY_NAME}"

# 设置配置文件路径
CONFIG_FILE="${PROJECT_DIR}/config.json"

# 创建构建目录和测试目录（如果不存在的话）
mkdir -p "$BUILD_DIR"
mkdir -p "$TEST_DIR"

# 1. 清理旧的构建文件
echo "Cleaning previous build..."
rm -f "$BINARY_PATH"

# 2. 构建 Go 项目
echo "Building the project..."
go build -o "$BINARY_PATH" .

# 检查构建是否成功
if [ $? -ne 0 ]; then
  echo "Build failed. Exiting..."
  exit 1
fi

echo "Build successful. Binary file located at: $BINARY_PATH"

# 3. 移动二进制文件到测试目录
echo "Moving binary file to test directory..."
cp "$BINARY_PATH" "$TEST_DIR"

# 4. 模拟生产环境（如果需要，你可以设置一些环境变量）
echo "Setting up production environment..."
export CONFIG_PATH="$CONFIG_FILE"  # 设置配置文件路径，供程序使用

# 5. 运行程序模拟生产环境
echo "Running the program in production mode..."
cd "$TEST_DIR" || exit
./"$BINARY_NAME"  # 运行构建好的二进制文件

# 6. 测试完后清理（如果需要）
echo "Cleaning up after test..."
rm -f "$TEST_DIR/$BINARY_NAME"

echo "Test completed successfully."
