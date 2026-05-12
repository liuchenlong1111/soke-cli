# 选择已有应用功能

## 功能概述

在用户登录授权流程中，新增了"选择已有应用"功能，允许用户从已创建的应用列表中选择一个应用进行配置，而不是每次都创建新应用。

## 实现内容

### 1. 新增页面

#### 选择应用页面 (`selectAppHTML`)
- 路由：`/select-app`
- 功能：展示用户企业下的所有应用列表
- 特性：
  - 应用卡片展示（包含应用名称、描述、图标）
  - 点击选择应用
  - 支持返回创建应用页面
  - 确认选择后跳转到成功页面

### 2. 新增 API 接口

#### 获取应用列表
- 路由：`GET /application/list`
- 功能：从授客开放平台获取用户企业下的应用列表
- 请求参数：`source=cli`（添加到 URL 查询参数）
- 后端 API：`GET {OpenDevURL}/app/admin/index?source=cli`
- 返回格式：
```json
{
  "success": true,
  "data": [
    {
      "app_id": "应用ID",
      "app_secret": "应用密钥",
      "name": "应用名称",
      "description": "应用描述"
    }
  ]
}
```

#### 后端 API 响应格式（授客开放平台）
```json
{
  "code": "200",
  "status": "ok",
  "message": "操作成功",
  "data": [
    {
      "id": 1,
      "app_id": "1D13B009-0482-4B80-B3C3-08CCE0C29E85",
      "apply_id": "31A57F40-0A93-44BA-B08E-B995333C1A01",
      "type": "corporation",
      "name": "soke测试",
      "description": "这是个测试申请",
      "app_key": "soke426c576c4ce58c02e80f5f2de28b65bd",
      "app_secret": "2189595a7a0cb1fee0c25f2b0fdc0760cbcc1662",
      "sso_secret": "brt4365uj75415ahj0207dc29b942ncc5l98pcb8"
    }
  ]
}
```

#### 选择应用
- 路由：`POST /application/select`
- 功能：根据 `app_id` 和 `app_key` 调用 `/app/admin/read` 接口获取完整应用信息，然后保存配置
- 请求格式：
```json
{
  "app_id": "应用ID",
  "app_key": "应用Key"
}
```
- 处理流程：
  1. 接收前端传来的 `app_id` 和 `app_key`
  2. 调用授客开放平台 `/app/admin/read` 接口获取完整应用信息（包括 `app_secret`）
  3. 从响应中提取 `app_secret`
  4. 保存 `app_id` 和 `app_secret` 到配置文件
- 返回格式：
```json
{
  "success": true,
  "message": "选择成功"
}
```

#### 读取应用详情接口（内部调用）
- 接口：`POST {OpenDevURL}/app/admin/read`
- 请求参数：`app_id`（form-urlencoded）
- 响应格式：
```json
{
  "code": "200",
  "status": "ok",
  "message": "操作成功",
  "data": {
    "app_id": "应用ID",
    "app_key": "应用Key",
    "app_secret": "应用密钥"
  }
}
```

### 3. 修改内容

#### 创建应用页面 (`createAppHTML`)
- 新增"选择已有应用"按钮
- 点击按钮跳转到 `/select-app` 页面

#### 路由新增
- `/create-app`：创建应用页面的别名路由，方便从选择页面返回

## 用户流程

1. 用户完成 OAuth 授权后，进入创建应用页面
2. 用户可以选择：
   - **创建新应用**：填写应用名称和描述，点击"创建"按钮
   - **选择已有应用**：点击"选择已有应用"按钮
3. 如果选择"选择已有应用"：
   - 跳转到应用列表页面
   - 系统自动加载用户企业下的所有应用（带 `source=cli` 参数）
   - 用户点击选择一个应用
   - 点击"确认"按钮
   - 系统保存选择的应用配置
   - 跳转到成功页面
4. 配置完成，用户可以开始使用 CLI

## 技术实现

### 前端
- 使用原生 JavaScript 实现
- Fetch API 进行异步请求
- 响应式 CSS 布局
- 交互式应用卡片选择

### 后端
- Go HTTP 路由处理
- 调用授客开放平台 API（添加 `source=cli` 参数）
- 解析授客开放平台返回的数据格式（包含 `code`、`status`、`message`、`data` 字段）
- **安全获取应用密钥**：通过 `/app/admin/read` 接口获取 `app_secret`，而不是直接从前端接收
- 配置文件读写（`core.CliConfig`）
- Token 数据管理

## 文件修改

- `cmd/user_auth/oauth_helpers.go`：新增 `selectAppHTML` 常量，修改 `createAppHTML`
- `cmd/user_auth/oauth_provider.go`：新增三个路由处理函数，修改 `/application/list` 接口

## 关键修改点

1. **添加 source 参数**：在请求应用列表时添加 `source=cli` 查询参数
2. **数据格式适配**：根据授客开放平台的实际返回格式解析数据
   - 响应包含 `code`、`status`、`message`、`data` 字段
   - `data` 数组中包含完整的应用信息（`id`、`app_id`、`app_key`、`app_secret` 等）
3. **数据转换**：将后端返回的数据转换为前端需要的简化格式（包含 `app_id`、`app_key`、`name`、`description`）
4. **安全获取密钥**：
   - 前端只传递 `app_id` 和 `app_key`
   - 后端通过调用 `/app/admin/read` 接口获取完整的应用信息
   - 从接口响应中提取 `app_secret` 并保存到配置文件
   - 避免在前端和网络传输中暴露敏感的 `app_secret`

## 测试建议

1. 测试创建新应用流程是否正常
2. 测试选择已有应用流程：
   - 有应用时的正常选择
   - 无应用时的提示
   - 网络错误时的处理
3. 测试页面跳转逻辑
4. 测试配置保存是否正确
5. 验证 `source=cli` 参数是否正确传递

## 注意事项

1. 应用列表接口依赖授客开放平台的 `/app/admin/index?source=cli` 接口
2. 应用详情接口依赖授客开放平台的 `/app/admin/read` 接口
3. 需要确保用户已完成 OAuth 授权，UserToken 有效
4. 选择应用后会覆盖配置文件中的 `app_id` 和 `app_secret`
5. 后端 API 返回格式包含 `code` 和 `status` 两个状态字段，需要同时检查
6. **安全性**：`app_secret` 不会在前端和客户端-服务器通信中传输，只在服务器端通过 API 获取
