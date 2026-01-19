#!/bin/bash
# Linktrack Install Script
# 安装脚本

set -e

INSTALL_DIR="/usr/local/bin"
REPO_URL="https://github.com/yourname/linktrack"

echo "Installing Linktrack..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go first: https://go.dev/dl/"
    exit 1
fi

# 创建临时目录
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# 克隆或下载
if command -v git &> /dev/null; then
    git clone --depth 1 "$REPO_URL" linktrack
    cd linktrack
else
    echo "Error: git is required."
    exit 1
fi

# 构建
make build

# 安装
sudo make install

# 初始化配置目录
make init

# 清理
rm -rf "$TMP_DIR"

echo ""
echo "Installation complete!"
echo "Run 'linktrack --help' to get started."
