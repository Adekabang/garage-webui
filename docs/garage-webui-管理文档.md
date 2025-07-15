# Garage Web UI 项目管理文档

## 项目概述

**Garage Web UI** 是一个用于管理 [Garage](https://garagehq.deuxfleurs.fr/) 分布式对象存储服务的现代化 Web 管理界面。该项目提供了一个简洁、直观的图形化界面来管理 Garage 集群，是 Garage 官方命令行工具的重要补充。

### 🎯 项目定位

- **目标用户**: Garage 集群管理员和运维人员
- **核心价值**: 简化 Garage 集群的日常管理操作
- **技术栈**: TypeScript + React (前端) + Go (后端)

## 功能特性

### 🏥 集群监控与管理

#### 1. 健康状态监控

- **实时集群状态**: 显示集群整体健康状况（健康/降级/不可用）
- **节点监控**: 实时监控已知节点数、连接节点数、存储节点状态
- **分区状态**: 监控数据分区的健康状况和仲裁状态

#### 2. 集群布局管理

- **可视化布局**: 图形化显示集群节点分布和存储配置
- **节点配置**: 管理节点的区域、容量、标签等属性
- **布局变更**: 支持暂存、预览、应用和回滚布局变更
- **历史记录**: 查看集群布局的历史变更记录

### 🗄️ 存储桶管理

#### 1. 存储桶操作

- **桶列表**: 显示所有存储桶及其基本信息
- **桶详情**: 查看存储桶的详细统计、配置和权限信息
- **桶创建**: 支持创建全局别名和本地别名的存储桶
- **桶配置**: 更新存储桶的网站配置、配额设置等

#### 2. 对象浏览器

- **文件浏览**: 内置对象浏览器，支持文件夹结构浏览
- **文件操作**: 上传、下载、删除对象文件
- **分享功能**: 生成临时访问链接
- **批量操作**: 支持批量文件管理

### 🔑 访问控制管理

#### 1. 访问密钥管理

- **密钥列表**: 显示所有 API 访问密钥
- **密钥创建**: 创建新的 S3 兼容访问密钥
- **权限配置**: 设置密钥的全局权限（如创建存储桶）
- **过期管理**: 设置密钥的过期时间

#### 2. 权限分配

- **桶权限**: 为访问密钥分配对特定存储桶的权限
- **权限类型**: 支持读取、写入、所有者三种权限级别
- **权限撤销**: 灵活的权限授予和撤销机制

## 技术架构

### 🏗️ 整体架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Browser   │───▶│  Garage Web UI   │───▶│  Garage Cluster │
│    (前端界面)    │    │   (Go 后端服务)   │    │  (Admin API)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 📁 项目结构

```
garage-webui/
├── src/                    # React 前端源码
│   ├── pages/             # 页面组件
│   │   ├── home/         # 首页仪表板
│   │   ├── cluster/      # 集群管理
│   │   ├── buckets/      # 存储桶管理
│   │   └── keys/         # 访问密钥管理
│   ├── components/       # 可复用组件
│   ├── hooks/           # React Hooks
│   └── lib/             # 工具库
├── backend/              # Go 后端源码
│   ├── main.go          # 服务入口
│   ├── router/          # API 路由
│   ├── utils/           # 工具函数
│   └── schema/          # 数据结构
├── docs/                # 项目文档
└── misc/                # 截图等资源
```

### 🔌 后端服务架构

#### 核心模块

1. **配置管理** (`utils/garage.go`)

   - 自动读取 Garage 配置文件 (`garage.toml`)
   - 提取管理 API 端点、认证信息等
   - 支持环境变量覆盖配置

2. **API 代理** (`router/`)

   - 代理前端请求到 Garage Admin API
   - 处理认证和错误转换
   - 提供统一的 RESTful 接口

3. **会话管理** (`utils/session.go`)

   - 支持用户认证（可选）
   - 会话状态管理

4. **缓存机制** (`utils/cache.go`)
   - API 响应缓存
   - 减少对 Garage 集群的请求压力

## 部署方案

### 🐳 Docker 部署（推荐）

#### 1. 与 Garage 集群一起部署

```yaml
services:
  garage:
    image: dxflrs/garage:v1.0.1
    volumes:
      - ./garage.toml:/etc/garage.toml
      - ./meta:/var/lib/garage/meta
      - ./data:/var/lib/garage/data
    ports:
      - 3900:3900 # S3 API
      - 3901:3901 # RPC
      - 3902:3902 # S3 Web
      - 3903:3903 # Admin API

  webui:
    image: khairul169/garage-webui:latest
    volumes:
      - ./garage.toml:/etc/garage.toml:ro
    ports:
      - 3909:3909
    environment:
      API_BASE_URL: "http://garage:3903"
      S3_ENDPOINT_URL: "http://garage:3900"
```

#### 2. 独立部署

```bash
docker run -p 3909:3909 \
  -v ./garage.toml:/etc/garage.toml:ro \
  -e API_BASE_URL="http://garage-host:3903" \
  -e API_ADMIN_KEY="your-admin-token" \
  khairul169/garage-webui:latest
```

### 🖥️ 二进制部署

```bash
# 下载二进制文件
wget -O garage-webui https://github.com/khairul169/garage-webui/releases/download/1.0.9/garage-webui-v1.0.9-linux-amd64
chmod +x garage-webui

# 运行服务
CONFIG_PATH=./garage.toml ./garage-webui
```

### 🔧 SystemD 服务

```ini
[Unit]
Description=Garage Web UI
After=network.target

[Service]
Environment="PORT=3909"
Environment="CONFIG_PATH=/etc/garage.toml"
ExecStart=/usr/local/bin/garage-webui
Restart=always

[Install]
WantedBy=default.target
```

## 配置管理

### 📝 Garage 配置要求

Web UI 需要 Garage 集群启用 Admin API：

```toml
# garage.toml
[admin]
api_bind_addr = "[::]:3903"
admin_token = "your-secure-admin-token"
metrics_token = "your-metrics-token"  # 可选
```

### 🌍 环境变量配置

| 变量名            | 描述                  | 默认值             |
| ----------------- | --------------------- | ------------------ |
| `CONFIG_PATH`     | Garage 配置文件路径   | `/etc/garage.toml` |
| `API_BASE_URL`    | Garage Admin API 地址 | 从配置文件读取     |
| `API_ADMIN_KEY`   | Admin API 令牌        | 从配置文件读取     |
| `S3_ENDPOINT_URL` | S3 API 地址           | 从配置文件读取     |
| `S3_REGION`       | S3 区域               | `garage`           |
| `BASE_PATH`       | Web UI 基础路径       | `/`                |
| `PORT`            | 服务端口              | `3909`             |
| `HOST`            | 绑定地址              | `0.0.0.0`          |

### 🔐 认证配置

#### 启用 Web UI 认证

```bash
# 生成密码哈希
htpasswd -nbBC 10 "admin" "password"

# 设置环境变量
AUTH_USER_PASS="admin:$2y$10$DSTi9o..."
```

## 管理最佳实践

### 🚀 日常运维

#### 1. 集群健康监控

- **定期检查**: 通过首页仪表板监控集群状态
- **告警设置**: 配置监控系统对接 `/metrics` 端点
- **性能观察**: 关注存储节点连接状态和分区健康度

#### 2. 存储桶管理

- **命名规范**: 建立统一的存储桶命名规范
- **权限最小化**: 为访问密钥分配最小必要权限
- **配额管理**: 为重要业务设置适当的配额限制

#### 3. 访问控制

- **定期轮换**: 定期轮换 API 访问密钥
- **权限审计**: 定期审查存储桶权限分配
- **密钥管理**: 为不同用途创建专用访问密钥

### 🔧 故障排查

#### 1. 连接问题

```bash
# 检查 Admin API 可访问性
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://garage-host:3903/v2/GetClusterHealth

# 检查网络连通性
telnet garage-host 3903
```

#### 2. 配置问题

- 验证 `garage.toml` 配置正确性
- 确认 Admin API 端口已开放
- 检查防火墙和网络策略

#### 3. 性能优化

- 启用缓存机制减少 API 调用
- 使用反向代理（如 Nginx）提供 SSL 终止
- 监控资源使用情况

### 📊 监控集成

#### Prometheus 指标

Web UI 可以配置为监控 Garage 的 Prometheus 指标：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: "garage"
    static_configs:
      - targets: ["garage-host:3903"]
    metrics_path: /metrics
    bearer_token: "your-metrics-token"
```

#### 关键指标

- `garage_cluster_health`: 集群健康状态
- `garage_storage_usage`: 存储使用情况
- `garage_api_requests`: API 请求统计
- `garage_replication_status`: 数据复制状态

## 开发指南

### 🛠️ 开发环境搭建

```bash
# 克隆项目
git clone https://github.com/khairul169/garage-webui.git
cd garage-webui

# 安装前端依赖
pnpm install

# 安装后端依赖
cd backend && go mod download && cd ..

# 启动开发服务器
pnpm run dev
```

### 🔧 技术选型说明

- **前端**: React 18 + TypeScript + Tailwind CSS
- **状态管理**: React Query (TanStack Query)
- **路由**: React Router
- **UI 组件**: 自定义组件库
- **后端**: Go + Gin 框架
- **配置解析**: go-toml

### 📋 贡献指南

1. **代码规范**: 遵循项目的 ESLint 和 Go fmt 规范
2. **测试**: 新功能需要添加相应测试
3. **文档**: 更新相关文档和 API 说明
4. **兼容性**: 确保与最新版本 Garage 兼容

## 安全考虑

### 🔒 安全建议

1. **网络安全**

   - 在生产环境中使用 HTTPS
   - 限制 Admin API 的网络访问
   - 使用防火墙规则保护敏感端口

2. **认证安全**

   - 启用 Web UI 用户认证
   - 使用强密码和定期轮换
   - 考虑集成企业身份认证系统

3. **权限控制**
   - 遵循最小权限原则
   - 定期审计访问权限
   - 使用专用的管理员 Token

## 未来规划

### 🚀 功能路线图

- **高级监控**: 集成更多性能指标和告警功能
- **批量操作**: 支持批量管理存储桶和访问密钥
- **API 扩展**: 支持更多 Garage Admin API 功能
- **国际化**: 多语言支持
- **主题系统**: 可定制的 UI 主题

### 🔧 技术改进

- **缓存优化**: 更智能的缓存策略
- **实时更新**: WebSocket 支持实时状态更新
- **移动优化**: 改进移动端体验
- **性能提升**: 前端打包优化和懒加载

## Garage Admin API 使用情况

### 🔌 当前项目调用的 API 功能

基于代码分析，当前 Garage Web UI 项目调用了以下 Garage Admin API v1 功能：

#### 1. 集群管理 API

- **`GET /v1/health`** - 获取集群健康状态

  - 用于首页仪表板显示集群状态
  - 监控节点连接数、存储节点状态、分区健康度

- **`GET /v1/status`** - 获取集群详细状态
  - 用于集群管理页面显示节点详情
  - 展示集群拓扑和节点配置信息

#### 2. 集群布局管理 API

- **`GET /v1/layout`** - 获取集群布局配置

  - 显示当前集群布局和暂存变更
  - 查看节点角色、容量、区域分配

- **`POST /v1/layout`** - 更新集群布局

  - 添加新节点到集群
  - 修改节点配置（容量、区域、标签）
  - 移除节点（设置 remove: true）

- **`POST /v1/connect`** - 连接集群节点

  - 将新节点连接到集群
  - 建立节点间的 RPC 连接

- **`POST /v1/layout/apply`** - 应用布局变更

  - 将暂存的布局变更应用到集群
  - 触发数据重新分布

- **`POST /v1/layout/revert`** - 回滚布局变更
  - 清除暂存的布局变更
  - 恢复到上一个稳定状态

#### 3. 存储桶管理 API

- **`GET /v1/bucket?list`** - 列出所有存储桶

  - 获取集群中所有存储桶列表
  - 显示桶的基本信息和别名

- **`GET /v1/bucket?id={id}`** - 获取存储桶详细信息

  - 查看单个存储桶的完整配置
  - 包含权限、统计、配额等信息

- **`POST /v1/bucket`** - 创建新存储桶

  - 支持设置全局别名和本地别名
  - 配置初始权限和参数

- **`PUT /v1/bucket?id={id}`** - 更新存储桶配置

  - 修改存储桶的网站配置
  - 设置或更新配额限制

- **`DELETE /v1/bucket?id={id}`** - 删除存储桶
  - 删除空的存储桶（需要桶为空）

#### 4. 存储桶别名管理 API

- **`PUT /v1/bucket/alias/global`** - 添加全局别名

  - 为存储桶创建全局访问别名
  - 支持多个别名指向同一个桶

- **`DELETE /v1/bucket/alias/global`** - 删除全局别名
  - 移除存储桶的全局别名
  - 保持桶本身不受影响

#### 5. 权限管理 API

- **`POST /v1/bucket/allow`** - 授予存储桶权限

  - 为访问密钥分配桶的操作权限
  - 支持读取、写入、所有者权限

- **`POST /v1/bucket/deny`** - 撤销存储桶权限
  - 移除访问密钥对桶的权限
  - 灵活的权限控制机制

#### 6. 访问密钥管理 API

- **`GET /v1/key?list`** - 列出所有访问密钥

  - 获取集群中的所有 API 密钥
  - 显示密钥的基本信息

- **`POST /v1/key`** - 创建新的访问密钥

  - 生成新的 S3 兼容访问密钥
  - 设置密钥的初始权限

- **`POST /v1/key/import`** - 导入已有访问密钥

  - 用于迁移或恢复访问密钥
  - 导入外部生成的密钥

- **`DELETE /v1/key?id={id}`** - 删除访问密钥
  - 从集群中移除访问密钥
  - 立即撤销所有相关权限

### ## API 版本对比分析

### 📊 当前项目 vs 官方文档 API 差异

通过对比分析，发现当前项目使用的是 **Garage Admin API v1**，而官方最新文档推荐使用 **API v2**。以下是详细的差异对比：

#### 🔄 版本映射关系

| 功能类别         | 当前项目 (v1)            | 官方推荐 (v2)                          | 状态      |
| ---------------- | ------------------------ | -------------------------------------- | --------- |
| **集群健康状态** | `GET /v1/health`         | `GET /v2/GetClusterHealth`             | ⚠️ 需升级 |
| **集群状态**     | `GET /v1/status`         | `GET /v2/GetClusterStatus`             | ⚠️ 需升级 |
| **集群统计**     | ❌ 未使用                | `GET /v2/GetClusterStatistics`         | 🆕 新功能 |
| **连接节点**     | `POST /v1/connect`       | `POST /v2/ConnectClusterNodes`         | ⚠️ 需升级 |
| **获取布局**     | `GET /v1/layout`         | `GET /v2/GetClusterLayout`             | ⚠️ 需升级 |
| **更新布局**     | `POST /v1/layout`        | `POST /v2/UpdateClusterLayout`         | ⚠️ 需升级 |
| **应用布局**     | `POST /v1/layout/apply`  | `POST /v2/ApplyClusterLayout`          | ⚠️ 需升级 |
| **回滚布局**     | `POST /v1/layout/revert` | `POST /v2/RevertClusterLayout`         | ⚠️ 需升级 |
| **布局历史**     | ❌ 未使用                | `GET /v2/GetClusterLayoutHistory`      | 🆕 新功能 |
| **预览布局变更** | ❌ 未使用                | `POST /v2/PreviewClusterLayoutChanges` | 🆕 新功能 |

#### 📦 存储桶管理 API 对比

| 功能           | 当前项目 (v1)                    | 官方推荐 (v2)                       | 差异说明            |
| -------------- | -------------------------------- | ----------------------------------- | ------------------- |
| **列出存储桶** | `GET /v1/bucket?list`            | `GET /v2/ListBuckets`               | 参数格式不同        |
| **获取桶信息** | `GET /v1/bucket?id={id}`         | `GET /v2/GetBucketInfo`             | 支持更多查询方式    |
| **创建存储桶** | `POST /v1/bucket`                | `POST /v2/CreateBucket`             | v2 支持更多配置选项 |
| **更新存储桶** | `PUT /v1/bucket?id={id}`         | `POST /v2/UpdateBucket/{id}`        | HTTP 方法和路径不同 |
| **删除存储桶** | `DELETE /v1/bucket?id={id}`      | `POST /v2/DeleteBucket/{id}`        | HTTP 方法不同       |
| **添加别名**   | `PUT /v1/bucket/alias/global`    | `POST /v2/AddBucketAlias`           | 支持本地别名        |
| **删除别名**   | `DELETE /v1/bucket/alias/global` | `POST /v2/RemoveBucketAlias`        | 支持本地别名        |
| **清理上传**   | ❌ 未使用                        | `POST /v2/CleanupIncompleteUploads` | 🆕 新功能           |
| **检查对象**   | ❌ 未使用                        | `GET /v2/InspectObject`             | 🆕 新功能           |

#### 🔑 访问密钥管理 API 对比

| 功能             | 当前项目 (v1)            | 官方推荐 (v2)             | 差异说明        |
| ---------------- | ------------------------ | ------------------------- | --------------- |
| **列出密钥**     | `GET /v1/key?list`       | `GET /v2/ListKeys`        | 参数格式不同    |
| **获取密钥信息** | ❌ 未使用                | `GET /v2/GetKeyInfo`      | 🆕 新功能       |
| **创建密钥**     | `POST /v1/key`           | `POST /v2/CreateKey`      | v2 支持更多选项 |
| **更新密钥**     | ❌ 未使用                | `POST /v2/UpdateKey/{id}` | 🆕 新功能       |
| **删除密钥**     | `DELETE /v1/key?id={id}` | `POST /v2/DeleteKey/{id}` | HTTP 方法不同   |
| **导入密钥**     | `POST /v1/key/import`    | `POST /v2/ImportKey`      | 路径结构不同    |
| **授予权限**     | `POST /v1/bucket/allow`  | `POST /v2/AllowBucketKey` | 路径结构不同    |
| **撤销权限**     | `POST /v1/bucket/deny`   | `POST /v2/DenyBucketKey`  | 路径结构不同    |

### 🚫 v2 独有功能（当前项目未使用）

#### 1. 管理员 Token 管理

- `GET /v2/ListAdminTokens` - 列出所有管理员 Token
- `GET /v2/GetAdminTokenInfo` - 获取 Token 信息
- `GET /v2/GetCurrentAdminTokenInfo` - 获取当前 Token 信息
- `POST /v2/CreateAdminToken` - 创建管理员 Token
- `POST /v2/UpdateAdminToken/{id}` - 更新管理员 Token
- `POST /v2/DeleteAdminToken/{id}` - 删除管理员 Token

#### 2. 节点管理

- `GET /v2/GetNodeInfo/{node}` - 获取节点信息
- `GET /v2/GetNodeStatistics/{node}` - 获取节点统计
- `POST /v2/CreateMetadataSnapshot/{node}` - 创建元数据快照
- `POST /v2/LaunchRepairOperation/{node}` - 启动修复操作

#### 3. 后台工作进程管理

- `POST /v2/ListWorkers/{node}` - 列出工作进程
- `POST /v2/GetWorkerInfo/{node}` - 获取工作进程信息
- `POST /v2/GetWorkerVariable/{node}` - 获取工作进程变量
- `POST /v2/SetWorkerVariable/{node}` - 设置工作进程变量

#### 4. 数据块管理

- `POST /v2/GetBlockInfo/{node}` - 获取数据块信息
- `GET /v2/ListBlockErrors/{node}` - 列出错误数据块
- `POST /v2/RetryBlockResync/{node}` - 重试数据块同步
- `POST /v2/PurgeBlocks/{node}` - 清除数据块

#### 5. 特殊端点

- `GET /health` - 快速健康检查（无需认证）
- `GET /metrics` - Prometheus 指标
- `GET /check` - 按需 TLS 检查

### ⚡ 升级影响分析

#### 🔴 关键差异

1. **API 路径结构**

   - v1: 使用查询参数 (`?id=xxx`)
   - v2: 使用 RESTful 路径 (`/{id}`)

2. **HTTP 方法**

   - v1: 混合使用 GET/POST/PUT/DELETE
   - v2: 主要使用 GET/POST

3. **请求/响应格式**

   - v2 提供更结构化的数据格式
   - 更详细的错误信息和状态码

4. **功能完整性**
   - v2 提供更多高级管理功能
   - 更好的监控和维护能力

#### 🟡 兼容性考虑

- **向后兼容**: v1 API 在当前版本中仍然可用（已标记为废弃）
- **迁移建议**: 逐步迁移到 v2 API
- **功能增强**: 利用 v2 新增功能改善用户体验

### 📋 升级建议

#### 🎯 短期计划（1-2 个月）

1. **API 版本升级**

   - 将核心 API 调用从 v1 升级到 v2
   - 更新前端 API 客户端
   - 测试兼容性和功能一致性

2. **基础功能增强**
   - 添加集群统计功能
   - 实现布局历史查看
   - 支持布局变更预览

#### 🚀 中期计划（3-6 个月）

1. **新功能集成**

   - 管理员 Token 管理界面
   - 节点详细信息和统计
   - 对象检查和分析功能

2. **监控增强**
   - 集成 Prometheus 指标显示
   - 实时健康状态监控
   - 错误和告警系统

#### 🎨 长期计划（6 个月以上）

1. **高级管理功能**

   - 数据块管理和修复工具
   - 后台工作进程监控
   - 自动化维护任务

2. **用户体验优化**
   - 批量操作支持
   - 实时数据更新
   - 移动端适配改进

### 📊 功能覆盖率分析

| 功能分类       | v1 可用功能 | v2 总功能 | 当前使用 | 覆盖率 |
| -------------- | ----------- | --------- | -------- | ------ |
| **集群管理**   | 4           | 6         | 2        | 33%    |
| **布局管理**   | 5           | 7         | 5        | 71%    |
| **存储桶管理** | 7           | 9         | 5        | 56%    |
| **权限管理**   | 2           | 2         | 2        | 100%   |
| **密钥管理**   | 4           | 6         | 4        | 67%    |
| **高级功能**   | 0           | 25+       | 0        | 0%     |
| **总体**       | 22          | 55+       | 18       | 33%    |

**结论**: 当前项目仅使用了 Garage Admin API 约 33% 的功能，有很大的功能扩展空间。
