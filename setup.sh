#!/bin/bash

# 设置错误时退出
set -e

# 设置数据库环境变量
export DB_USER="app_user"
export DB_PASSWORD="app_password"
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_NAME="app_db"

# 设置Redis环境变量
export REDIS_HOST="localhost"
export REDIS_PORT="6379"

# 检查是否安装了Homebrew
if ! command -v brew &> /dev/null; then
    echo "Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# 安装Go
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    brew install go
fi

# 安装Docker
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    brew install --cask docker
    echo "Please open Docker.app to complete the installation"
fi

# 安装Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    brew install docker-compose
fi

# 安装pre-commit
if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."
    brew install pre-commit
fi

# 安装golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    brew install golangci-lint
fi

# 安装HTTPie
if ! command -v http &> /dev/null; then
    echo "Installing HTTPie..."
    brew install httpie
fi

# 初始化pre-commit钩子
echo "Initializing pre-commit hooks..."
pre-commit install

# 验证安装
echo "Verifying installations..."
echo "Go version: $(go version)"
echo "Docker version: $(docker --version)"
echo "Docker Compose version: $(docker-compose --version)"
echo "pre-commit version: $(pre-commit --version)"
echo "golangci-lint version: $(golangci-lint --version)"
echo "HTTPie version: $(http --version)"

# 设置执行权限
chmod +x setup.sh

echo "Setup completed successfully!"
echo "You can now run 'make dev' to start the development environment."
