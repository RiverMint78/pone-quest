# PoneQuest

![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-v4-06B6D4?logo=tailwindcss&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-2.x-3366CC?logo=htmx&logoColor=white)

PoneQuest 是一个 MLP 台词数据的本地语义检索工具。
使用本地 embedding 模型和本地向量索引完成查询，不依赖云端检索服务。

灵感来源：

- zvv.quest (MemeMeow)：<https://github.com/MemeMeow-Studio/MemeMeow>
- poneponepone.org (Transcript Search)：<https://github.com/ReluctusB/MLP-Fim-Episode-Transcript-Search>

## 项目目标

- 用自然语言搜索 MLP 台词，而不只靠关键词精确匹配。
- 保持本地运行，方便个人部署。
- 简单稳定的 Web 界面和后端接口。

## 主要功能

- 本地 embedding + 本地索引，无云检索，无数据库依赖
- Go 后端 + Bun 前端构建
- 支持台词检索、剧集详情与健康检查

## 工作流程

1. 准备台词数据（`episodes.json`）。
2. 使用 `cmd/indexer` 生成或更新向量索引（`lines.idx`）。
3. 启动 `cmd/server`，加载模型、数据和索引。
4. 通过浏览器访问站点并发起检索。

## 技术栈

- Go 1.26
- Bun (前端构建)
- TypeScript + Tailwind CSS + HTMX
- GGUF 本地 embedding 模型
- [kelindar/search](https://github.com/kelindar/search) 向量索引

## 环境要求

- Go 1.26+
- Bun 1.x
- 可用的 GGUF embedding 模型文件
- [kelindar/search](https://github.com/kelindar/search)  所需的轻量 llama.cpp 绑定库

## 快速开始

### 1) 安装前端依赖

```bash
bun install
```

### 2) 构建前端静态资源

```bash
bun run build
```

### 3) 配置 `.env`

在项目根目录创建 `.env` 文件：

```env
PQ_EMBEDDING_MODEL=data/models/bge-base-zh-v1.5-q4_k_m.gguf
PQ_EPISODEITEM_FILE=data/episodes.json
PQ_INDEX_FILE=data/lines.idx
PQ_QUERY_INSTRUCTION=
```

### 4) 生成索引

```bash
go run ./cmd/indexer
```

你需要准备一个类似于 `data/episodes.json` 的台本文件。

### 5) 启动服务

```bash
# 如果要展示详细调试信息
go run ./cmd/server -port 8080 -debug
```

访问：<http://localhost:8080>

## 开发命令

前端：

```bash
bun run dev
bun run build
```

后端：

```bash
go test ./... # 然而测试只不过随便写点
go vet ./...
```

## 环境变量

| 变量 | 必填 | 说明 |
| --- | --- | --- |
| `PQ_EMBEDDING_MODEL` | 是 | 本地模型路径 |
| `PQ_EPISODEITEM_FILE` | 是 | 台本数据 JSON 路径 |
| `PQ_INDEX_FILE` | 是 | 向量索引文件路径 |
| `PQ_QUERY_INSTRUCTION` | 否 | 查询前缀 |

## HTTP 端点

- `GET /`：主页
- `GET /search`：搜索
- `GET /episode`：剧集明细
- `GET /healthz`：健康检查
- `GET /static/*`：静态资源
