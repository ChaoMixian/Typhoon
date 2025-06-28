Typhoon API 手册 (初稿)

版本: 1.0.0 基础URL: http://<your_typhoon_host>:<port>/api/v1 (请替换为您的实际部署地址和端口)

认证: 所有需要认证的 API 请求都应在 HTTP Header 中包含 Authorization 字段，值为 Bearer <your_api_token>。 <your_api_token> 在 Typhoon 的 config.json 中设置。

1. 配置管理 (Typhoon 自身配置)

这些端点用于管理 Typhoon 应用本身的配置 (config.json)。

1.1 获取当前配置

GET /config
描述: 获取 Typhoon 当前的完整配置信息。
响应:
200 OK: 成功，返回配置对象的 JSON 结构。
500 Internal Server Error: 获取配置失败。
示例响应 (application/json):
{
  "proxy": {
    "currentCore": "mihomo",
    "mihomo": {
      "binPath": "/path/to/mihomo",
      "currentConfig": "default",
      "controllerAddress": "0.0.0.0:9999",
      "listenPort": 7890,
      "tun": {
        "enabled": false,
        "stack": "system",
        "dnsHijaking": "0.0.0.0:53"
      }
    },
    "mode": "transparent",
    "dns": { /* DNS 配置详情 */ }
  },
  "logging": { /* 日志配置详情 */ },
  "subscriptionManage": { /* 订阅管理配置详情 */ },
  "api": { /* API 服务配置详情 */ }
}
1.2 更新部分配置

PATCH /config/update
描述: 更新 Typhoon 配置中的一个或多个部分。请求体应包含要更新的配置项。
请求体 (application/json):
{
  "proxy": { // 可选
    "currentCore": "xray", // 可选
    "mode": "rule", // 可选
    "mihomo": { "listenPort": 7891 }, // 可选, 只更新指定字段
    "dns": { "enabled": false, "listen": "0.0.0.0:5353" } // 可选
  },
  "logging": { "level": "debug" }, // 可选
  "api": { "token": "newSecretToken" } // 可选
}
注意: 对于嵌套对象 (如 mihomo, dns), 如果提供了该键，则其内部所有字段都将被视为该部分的完整新状态（除非后端逻辑支持更深层次的合并）。对于简单字段（如 currentCore, mode, level），则直接更新。
响应:
200 OK: 配置更新成功。
400 Bad Request: 请求体格式错误或无效的配置项。
500 Internal Server Error: 更新配置失败。
示例响应 (application/json):
{
  "message": "Configuration updated successfully"
}
1.3 重新加载配置

POST /config/reload
描述: 强制 Typhoon 从磁盘重新加载 config.json 文件。
响应:
200 OK: 配置重新加载成功。
500 Internal Server Error: 重新加载失败（例如，配置文件格式错误）。服务器可能会继续使用旧的有效配置。
示例响应 (application/json):
{
  "message": "Configuration reloaded successfully"
}
2. 守护进程管理 (Mihomo 核心控制)

这些端点用于控制 Mihomo 核心进程。

2.1 启动 Mihomo 核心

POST /daemon/mihomo/start
描述: 根据 Typhoon 当前配置启动 Mihomo 核心进程。
响应:
200 OK: Mihomo 启动成功。
400 Bad Request: Mihomo 已在运行或配置错误（如找不到二进制文件）。
500 Internal Server Error: 启动过程中发生内部错误。
示例响应 (application/json):
{
  "message": "Mihomo daemon started successfully"
}
2.2 停止 Mihomo 核心

POST /daemon/mihomo/stop
描述: 停止当前运行的 Mihomo 核心进程。
响应:
200 OK: Mihomo 停止成功。
400 Bad Request: Mihomo 未在运行。
500 Internal Server Error: 停止过程中发生内部错误。
示例响应 (application/json):
{
  "message": "Mihomo daemon stopped successfully"
}
2.3 重启 Mihomo 核心

POST /daemon/mihomo/restart
描述: 重启 Mihomo 核心进程（先停止后根据当前配置启动）。
响应:
200 OK: Mihomo 重启成功。
500 Internal Server Error: 重启过程中发生内部错误。
示例响应 (application/json):
{
  "message": "Mihomo daemon restarted successfully"
}
2.4 获取 Mihomo 核心状态

GET /daemon/mihomo/status
描述: 获取 Mihomo 核心的当前运行状态。
响应 (application/json):
200 OK:
{
  "isRunning": true, // 布尔值，表示是否正在运行
  "pid": 12345,    // 如果正在运行，则为进程ID
  "version": "v1.x.x" // Mihomo 核心版本 (如果可获取)
}
3. 订阅管理

3.1 更新所有订阅

POST /subscriptions/update
描述: 触发所有已配置的订阅链接进行更新。
响应:
200 OK: 订阅更新任务已成功处理（部分或全部成功）。
{
  "status": "Subscriptions updated successfully",
  "results": [
    { "name": "Subscription A", "success": true },
    { "name": "Subscription B", "success": false, "error": "download failed" }
  ]
}
500 Internal Server Error: 更新过程中发生严重错误。
4. Mihomo API 代理

Typhoon 代理了 Mihomo 的原生 API。您可以通过 Typhoon 访问这些 API，路径通常是在 Typhoon 的 API 基础上加上一个前缀，例如 /mihomo/。

代理基础路径: /api/v1/mihomo (假设)

这意味着，如果 Mihomo 的一个端点是 /logs，通过 Typhoon 访问它将是 /api/v1/mihomo/logs。 所有 Mihomo 原生 API 的认证（如果 Mihomo 本身配置了 secret）和参数都应按 Mihomo 文档进行传递。 Typhoon 的 Authorization: Bearer <typhoon_token> 仍然需要，用于访问 Typhoon 的代理本身。

4.1 获取 Mihomo 日志 (示例)

GET /mihomo/logs
描述: 获取 Mihomo 核心的实时日志 (流式数据)。
(具体请求和响应格式遵循 Mihomo API 文档)
4.2 获取 Mihomo 流量信息 (示例)

GET /mihomo/traffic
描述: 获取 Mihomo 核心的实时流量信息。
(具体请求和响应格式遵循 Mihomo API 文档)
4.3 获取 Mihomo 版本信息 (示例)

GET /mihomo/version
描述: 获取 Mihomo 核心的版本。
(具体请求和响应格式遵循 Mihomo API 文档)
4.4 更新 Mihomo 运行配置 (示例)

PATCH /mihomo/configs
描述: 更新 Mihomo 核心的部分运行配置。
请求体 (application/json): (遵循 Mihomo API 文档，例如 {"mixed-port": 7890})
(具体请求和响应格式遵循 Mihomo API 文档)
4.5 其他 Mihomo API 端点

其他如策略组管理 (/group/...)、代理节点选择 (/proxies/...)、连接管理 (/connections/...) 等均可通过此代理机制访问。请参考您之前提供的 Mihomo API 文档，并将路径前缀替换为 /api/v1/mihomo。
错误处理:

400 Bad Request: 通常表示客户端请求参数错误、格式无效或缺少必要字段。
401 Unauthorized: API Token 未提供或无效。
403 Forbidden: API Token 有效，但无权执行该操作。
404 Not Found: 请求的资源或路径不存在。
500 Internal Server Error: 服务器内部发生错误。
响应体中通常会包含一个 error 字段提供更详细的错误信息，例如:

{
  "error": "Descriptive error message here"
}
这只是一个初步的框架。您需要根据后端的具体实现来填充和校对：

确切的路由: 确保所有路由都是您后端 Gin 中定义的路由。
请求/响应体: 详细列出每个端点预期的请求体字段和确切的响应体结构，包括成功和失败的示例。
Mihomo 代理前缀: 确认 Mihomo API 代理所使用的确切路径前缀 (我假设了 /mihomo/)。
更细致的配置项: 在 /config (PATCH) 中，明确哪些深层嵌套的配置项支持部分更新，哪些需要提供完整对象。
SSE (Server-Sent Events): 对于像日志流或更新进度这样的端点，需要特别注明它们使用的是 SSE。
