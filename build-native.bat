@echo off
REM Windows 原生镜像构建脚本

echo 开始清理项目...
call mvn clean

echo 开始构建原生镜像...
call mvn package -Pnative -e -X

echo 构建完成！
pause