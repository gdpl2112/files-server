#!/bin/bash
# 原生镜像构建脚本

echo "开始清理项目..."
mvn clean

echo "开始构建原生镜像..."
mvn package -Pnative -e -X

echo "构建完成！"