# OpenAI API 代理服务器

这是一个简单的代理服务器，用于将无需 API Key 的 OpenAI 格式 API 转换为需要 API Key 的格式。

## 功能

- 将无需 API Key 的请求转发到目标 API
- 自动添加模拟的 API Key
- 支持所有 OpenAI API 端点

## 安装

确保你已安装 Go 1.21 或更高版本，然后运行：

```bash
go mod download
```

## 使用方法

1. 编译程序：
```bash
go build -o proxy_server
```

2. 运行服务器：
```bash
./proxy_server -target="https://你的目标API地址" -key="sk-你想使用的API-Key" -addr=":8080"
```

参数说明：
- `-target`: 必需，指定目标 API 的基础 URL
- `-key`: 可选，指定要使用的 API Key，默认为 "sk-dummy-key"
- `-addr`: 可选，指定服务器监听地址，默认为 ":8080"

3. 在 Cherry Studio 中配置：
   - API 地址：`http://localhost:8080`
   - API Key：使用你在启动代理服务器时指定的 key

## 示例

```bash
./proxy_server -target="https://api.example.com" -key="sk-your-key" -addr=":8080"
```

这将启动一个代理服务器，监听在 8080 端口，将所有请求转发到 https://api.example.com，并添加指定的 API Key。 