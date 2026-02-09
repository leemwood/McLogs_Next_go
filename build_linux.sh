#!/bin/bash

# Linux 编译脚本
BUILD_DIR="dist/linux"
mkdir -p $BUILD_DIR/configs

echo "Building Go backend for Linux..."
# 交叉编译环境设置
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/mclogs-api ./cmd/server/main.go

echo "Copying configurations..."
cp configs/config.yaml $BUILD_DIR/configs/config.yaml
cp configs/patterns.yaml $BUILD_DIR/configs/patterns.yaml

echo "Linux build complete in $BUILD_DIR"
