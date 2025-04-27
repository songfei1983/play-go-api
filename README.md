# Play Go API

一个基于Go语言的现代化API服务示例项目，集成了多种云原生技术和最佳实践。

## 项目概述

本项目是一个完整的Go语言API服务示例，展示了如何构建一个具有以下特性的现代化微服务：

- RESTful API设计
- 数据库集成与缓存
- 分布式追踪
- 监控与指标收集
- API网关
- 容器化部署

## 技术栈

### 后端

- **Go**: 主要编程语言
- **Echo**: 高性能、可扩展的Web框架
- **GORM**: ORM库，用于数据库操作
- **Redis**: 用于缓存

### 数据存储

- **MySQL**: 主数据库
- **Redis**: 缓存层

### 可观测性

- **OpenTelemetry**: 分布式追踪框架
- **Jaeger**: 分布式追踪系统
- **Prometheus**: 监控和时间序列数据库
- **Grafana**: 可视化和监控平台

### API网关

- **Kong**: 高性能API网关

### 容器化

- **Docker**: 容器化应用
- **Docker Compose**: 多容器应用编排

## 系统架构

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│    Client   │──────▶    Kong     │──────▶   Go API     │
└─────────────┘      │  API Gateway │      │   Service   │
                     └─────────────┘      └──────┬──────┘
                                                  │
                                                  ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Grafana   │◀─────│ Prometheus  │◀─────│   Metrics   │
└─────────────┘      └─────────────┘      └─────────────┘
      ▲                                          ▲
      │                                          │
      │            ┌─────────────┐              │
      └────────────│   Jaeger    │◀─────────────┘
                   └─────────────┘
                         ▲
                         │
┌─────────────┐      ┌─────────────┐
│    Redis    │◀─────▶    MySQL    │
└─────────────┘      └─────────────┘
```

## 功能特性

- **用户管理**: 注册和获取用户信息
- **缓存层**: 使用Redis缓存用户数据
- **分布式追踪**: 使用OpenTelemetry和Jaeger追踪请求流程
- **监控**: 使用Prometheus收集指标，Grafana展示仪表盘
- **API网关**: 使用Kong进行路由和负载均衡
- **优雅关闭**: 支持服务的优雅启动和关闭

## 安装与运行

### 前提条件

- Go 1.16+
- Docker 和 Docker Compose
- Make

### 本地开发

1. 克隆仓库

```bash
git clone https://github.com/songfei1983/play-go-api.git
cd play-go-api
```

2. 启动依赖服务

```bash
make docker-up
```

3. 运行应用

```bash
make run
```

或者使用开发模式：

```bash
make dev
```

## Local Development

### Prerequisites
- Go 1.21 or higher
- MySQL
- Redis

### Setting up local environment
1. Copy the environment example file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your local configuration if needed.

3. Run the application locally:
```bash
make run-local
```

This will:
- Set up the required environment variables
- Build the application
- Start the server locally

You can also set environment variables manually using:
```bash
make set-env
```

### 使用Docker Compose

一键启动所有服务：

```bash
make docker-build
make docker-up
```

停止所有服务：

```bash
make docker-down
```

## API文档

启动服务后，可以通过以下方式访问API文档：

- Swagger UI: http://localhost:8080/swagger/index.html

## 监控与追踪

- Grafana: http://localhost:3000 (用户名: admin, 密码: admin)
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090

## 项目结构

```
.
├── cmd/                # 应用入口
│   └── api/            # API服务入口
├── config/             # 配置文件
│   ├── grafana/        # Grafana配置和仪表盘
│   └── ...             # 其他配置文件
├── docs/               # API文档
├── internal/           # 内部包
│   ├── app/            # 应用核心
│   ├── config/         # 配置加载
│   ├── handler/        # HTTP处理器
│   ├── metrics/        # 指标收集
│   ├── middleware/     # HTTP中间件
│   └── server/         # HTTP服务器
├── Dockerfile          # Docker构建文件
├── docker-compose.yml  # Docker Compose配置
├── Makefile            # 构建和运行命令
└── README.md           # 项目文档
```

## 环境变量

应用支持通过环境变量进行配置：

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| DB_HOST | MySQL主机 | - |
| DB_PORT | MySQL端口 | - |
| DB_USER | MySQL用户名 | - |
| DB_PASSWORD | MySQL密码 | - |
| DB_NAME | MySQL数据库名 | - |
| REDIS_HOST | Redis主机 | - |
| REDIS_PORT | Redis端口 | - |
| SERVER_PORT | API服务端口 | 8080 |
| TRACING_ENDPOINT | Jaeger端点 | jaeger:4317 |

## 贡献

欢迎提交问题和拉取请求！

## 许可证

[MIT](LICENSE)
