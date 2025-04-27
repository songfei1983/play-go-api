#!/bin/bash

# 确保 API 服务正在运行
echo "确保 API 服务和依赖服务正在运行..."
cd /Users/songfei/Develop/github.com/songfei1983/play-go-api
docker-compose up -d

# 等待服务就绪
echo "等待服务就绪..."
sleep 10

# 创建结果目录
mkdir -p test/gatling/results

# 构建并运行 Gatling 测试容器
cd test/gatling
docker-compose build
docker-compose run --rm gatling

# 输出测试报告位置
echo "测试完成！结果报告在 test/gatling/results 目录中"

# 打开最新的测试报告
latest_report=$(ls -t results/*/index.html | head -n 1)
if [ -n "$latest_report" ]; then
    echo "打开测试报告: $latest_report"
    open "$latest_report"
fi
