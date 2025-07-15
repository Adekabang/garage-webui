# Garage Admin API 文档

## 概述

Garage Administration API 是一个用于编程式管理 Garage 集群的 REST API，提供了完整的集群管理、存储桶管理、访问控制等功能。当前版本为 v2，API 基础地址通常为 `http://localhost:3903`。

## 认证方式

### Bearer Token 认证

所有 API 请求都需要在 HTTP 头中包含认证信息：

```http
Authorization: Bearer <token>
```

### Token 类型

1. **用户定义 Token**（推荐）

   - 可动态创建和管理
   - 支持作用域限制
   - 支持过期时间设置
   - 使用 `garage admin-token` 命令创建

2. **主 Token**（已废弃）
   - 在配置文件中指定
   - `admin_token`: 管理端点访问
   - `metrics_token`: 指标端点访问

### 创建用户定义 Token 示例

```bash
garage admin-token create --expires-in 30d \
    --scope ListBuckets,GetBucketInfo,ListKeys,GetKeyInfo,CreateBucket,CreateKey,AllowBucketKey,DenyBucketKey \
    my-token
```

## API 端点分类

### 1. 集群管理 (Cluster)

#### 获取集群健康状态

- **端点**: `GET /v2/GetClusterHealth`
- **描述**: 返回集群全局状态，包括连接节点数、健康存储节点数、分区状态等
- **响应示例**:

```json
{
  "status": "healthy",
  "knownNodes": 3,
  "connectedNodes": 3,
  "storageNodes": 3,
  "storageNodesOk": 3,
  "partitions": 256,
  "partitionsQuorum": 256,
  "partitionsAllOk": 256
}
```

#### 获取集群状态

- **端点**: `GET /v2/GetClusterStatus`
- **描述**: 返回详细的集群状态信息，包括节点信息和布局配置

#### 获取集群统计

- **端点**: `GET /v2/GetClusterStatistics`
- **描述**: 获取集群级别的统计数据

#### 连接集群节点

- **端点**: `POST /v2/ConnectClusterNodes`
- **描述**: 指示当前节点连接到其他 Garage 节点
- **请求体**: 节点地址数组 `["<node_id>@<net_address>"]`

### 2. 集群布局管理 (Cluster Layout)

#### 获取集群布局

- **端点**: `GET /v2/GetClusterLayout`
- **描述**: 返回当前集群布局配置和待处理的变更

#### 更新集群布局

- **端点**: `POST /v2/UpdateClusterLayout`
- **描述**: 提交集群布局变更到暂存区
- **请求体示例**:

```json
{
  "roles": [
    {
      "id": "node-id",
      "zone": "zone1",
      "capacity": 100000000000,
      "tags": ["tag1", "tag2"]
    }
  ]
}
```

#### 应用布局变更

- **端点**: `POST /v2/ApplyClusterLayout`
- **描述**: 将暂存的布局变更应用到集群
- **请求体**: `{"version": <layout_version>}`

#### 预览布局变更

- **端点**: `POST /v2/PreviewClusterLayoutChanges`
- **描述**: 预览布局变更的影响，不实际应用

#### 回滚布局变更

- **端点**: `POST /v2/RevertClusterLayout`
- **描述**: 清除所有暂存的布局变更

#### 获取布局历史

- **端点**: `GET /v2/GetClusterLayoutHistory`
- **描述**: 获取集群布局的历史版本信息

### 3. 存储桶管理 (Bucket)

#### 列出所有存储桶

- **端点**: `GET /v2/ListBuckets`
- **描述**: 返回集群中所有存储桶及其别名

#### 获取存储桶信息

- **端点**: `GET /v2/GetBucketInfo`
- **参数**:
  - `id`: 存储桶 ID
  - `globalAlias`: 全局别名
  - `search`: 搜索模式
- **描述**: 获取存储桶详细信息，包括权限、统计、配额等

#### 创建存储桶

- **端点**: `POST /v2/CreateBucket`
- **请求体示例**:

```json
{
  "globalAlias": "my-bucket",
  "localAlias": {
    "accessKeyId": "key-id",
    "alias": "local-name",
    "allow": {
      "read": true,
      "write": true,
      "owner": false
    }
  }
}
```

#### 更新存储桶

- **端点**: `POST /v2/UpdateBucket/{id}`
- **描述**: 更新存储桶的网站配置和配额设置
- **请求体示例**:

```json
{
  "websiteAccess": {
    "enabled": true,
    "indexDocument": "index.html",
    "errorDocument": "error.html"
  },
  "quotas": {
    "maxSize": 1000000000,
    "maxObjects": 10000
  }
}
```

#### 删除存储桶

- **端点**: `POST /v2/DeleteBucket/{id}`
- **描述**: 删除空存储桶（会删除所有关联别名）

#### 清理未完成上传

- **端点**: `POST /v2/CleanupIncompleteUploads`
- **请求体**: `{"bucketId": "bucket-id", "olderThanSecs": 86400}`

#### 检查对象

- **端点**: `GET /v2/InspectObject`
- **参数**: `bucketId`, `key`
- **描述**: 获取对象的详细内部状态信息

### 4. 存储桶别名管理 (Bucket Alias)

#### 添加存储桶别名

- **端点**: `POST /v2/AddBucketAlias`
- **描述**: 为存储桶添加全局或本地别名

#### 移除存储桶别名

- **端点**: `POST /v2/RemoveBucketAlias`
- **描述**: 移除存储桶的别名

### 5. 访问密钥管理 (Access Key)

#### 列出访问密钥

- **端点**: `GET /v2/ListKeys`
- **描述**: 返回所有 API 访问密钥

#### 获取密钥信息

- **端点**: `GET /v2/GetKeyInfo`
- **参数**:
  - `id`: 密钥 ID
  - `search`: 搜索模式
  - `showSecretKey`: 是否返回密钥（默认不返回）

#### 创建访问密钥

- **端点**: `POST /v2/CreateKey`
- **请求体示例**:

```json
{
  "name": "my-key",
  "allow": {
    "createBucket": true
  }
}
```

#### 更新访问密钥

- **端点**: `POST /v2/UpdateKey/{id}`
- **描述**: 更新密钥的名称、权限和过期时间

#### 删除访问密钥

- **端点**: `POST /v2/DeleteKey/{id}`
- **描述**: 从集群中删除访问密钥

#### 导入访问密钥

- **端点**: `POST /v2/ImportKey`
- **描述**: 导入已有的访问密钥（仅用于迁移和备份恢复）

### 6. 权限管理 (Permission)

#### 授予权限

- **端点**: `POST /v2/AllowBucketKey`
- **描述**: 授予密钥对存储桶的操作权限
- **请求体示例**:

```json
{
  "bucketId": "bucket-id",
  "accessKeyId": "key-id",
  "permissions": {
    "read": true,
    "write": true,
    "owner": false
  }
}
```

#### 拒绝权限

- **端点**: `POST /v2/DenyBucketKey`
- **描述**: 移除密钥对存储桶的操作权限

### 7. 管理员 Token 管理 (Admin API Token)

#### 列出管理员 Token

- **端点**: `GET /v2/ListAdminTokens`

#### 获取 Token 信息

- **端点**: `GET /v2/GetAdminTokenInfo`
- **参数**: `id` 或 `search`

#### 获取当前 Token 信息

- **端点**: `GET /v2/GetCurrentAdminTokenInfo`

#### 创建管理员 Token

- **端点**: `POST /v2/CreateAdminToken`
- **请求体示例**:

```json
{
  "name": "my-admin-token",
  "expiration": "2025-12-31T23:59:59Z",
  "scope": ["ListBuckets", "GetBucketInfo", "CreateBucket"]
}
```

#### 更新管理员 Token

- **端点**: `POST /v2/UpdateAdminToken/{id}`

#### 删除管理员 Token

- **端点**: `POST /v2/DeleteAdminToken/{id}`

### 8. 节点管理 (Node)

#### 获取节点信息

- **端点**: `GET /v2/GetNodeInfo/{node}`
- **参数**: `node` - 节点 ID、`*`（所有节点）或 `self`（当前节点）

#### 获取节点统计

- **端点**: `GET /v2/GetNodeStatistics/{node}`

#### 创建元数据快照

- **端点**: `POST /v2/CreateMetadataSnapshot/{node}`

#### 启动修复操作

- **端点**: `POST /v2/LaunchRepairOperation/{node}`
- **修复类型**: `tables`, `blocks`, `versions`, `multipartUploads`, `blockRefs`, `blockRc`, `rebalance`, `aliases`

### 9. 后台工作进程管理 (Worker)

#### 列出工作进程

- **端点**: `POST /v2/ListWorkers/{node}`

#### 获取工作进程信息

- **端点**: `POST /v2/GetWorkerInfo/{node}`

#### 获取工作进程变量

- **端点**: `POST /v2/GetWorkerVariable/{node}`

#### 设置工作进程变量

- **端点**: `POST /v2/SetWorkerVariable/{node}`

### 10. 数据块管理 (Block)

#### 获取数据块信息

- **端点**: `POST /v2/GetBlockInfo/{node}`
- **请求体**: `{"blockHash": "hash-value"}`

#### 列出错误数据块

- **端点**: `GET /v2/ListBlockErrors/{node}`

#### 重试数据块同步

- **端点**: `POST /v2/RetryBlockResync/{node}`

#### 清除数据块

- **端点**: `POST /v2/PurgeBlocks/{node}`
- **警告**: 此操作会永久删除引用这些数据块的所有对象

### 11. 特殊端点 (Special Endpoints)

#### 健康检查

- **端点**: `GET /health`
- **认证**: 无需认证
- **描述**: 快速健康检查，返回 200 表示服务可用

#### Prometheus 指标

- **端点**: `GET /metrics`
- **认证**: 可选（使用 metrics_token）
- **描述**: 返回 Prometheus 格式的监控指标

#### 按需 TLS 检查

- **端点**: `GET /check?domain=<domain>`
- **认证**: 无需认证
- **描述**: 用于反向代理（如 Caddy）的按需 TLS 证书验证

## 使用示例

### 使用 curl

```bash
# 获取集群健康状态
curl -H 'Authorization: Bearer YOUR_TOKEN' \
     http://localhost:3903/v2/GetClusterHealth

# 创建存储桶
curl -X POST \
     -H 'Authorization: Bearer YOUR_TOKEN' \
     -H 'Content-Type: application/json' \
     -d '{"globalAlias": "my-bucket"}' \
     http://localhost:3903/v2/CreateBucket

# 列出所有存储桶
curl -H 'Authorization: Bearer YOUR_TOKEN' \
     http://localhost:3903/v2/ListBuckets
```

### 使用 Garage CLI

```bash
# 通过内部 RPC 调用（无需认证）
garage json-api GetClusterHealth

# 带参数的调用
garage json-api GetBucketInfo '{"globalAlias": "my-bucket"}'

# 从标准输入读取参数
garage json-api CreateBucket -
{"globalAlias": "test-bucket"}
<EOF>
```

## 错误处理

API 使用标准 HTTP 状态码：

- `200 OK` - 请求成功
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 认证失败
- `403 Forbidden` - 权限不足
- `404 Not Found` - 资源不存在
- `500 Internal Server Error` - 服务器内部错误

错误响应通常包含详细的错误信息：

```json
{
  "error": "Bucket not found",
  "code": "BucketNotFound"
}
```

## 权限作用域

管理员 Token 可以限制访问特定的 API 端点：

- `*` - 允许所有端点
- `ListBuckets` - 列出存储桶
- `GetBucketInfo` - 获取存储桶信息
- `CreateBucket` - 创建存储桶
- `ListKeys` - 列出访问密钥
- `CreateKey` - 创建访问密钥
- `AllowBucketKey` - 授予权限
- `DenyBucketKey` - 拒绝权限
- `Metrics` - 访问指标端点

## 最佳实践

1. **使用用户定义 Token**：避免使用配置文件中的主 Token
2. **设置适当的作用域**：只授予必要的权限
3. **设置过期时间**：定期轮换 Token
4. **监控 API 使用**：通过 `/metrics` 端点监控 API 调用
5. **错误处理**：妥善处理各种错误情况
6. **批量操作**：对于大量操作，考虑使用批量 API 或脚本

## 版本历史

- **v0** - Garage v0.7.2 首次引入（已废弃）
- **v1** - Garage v0.9.0 引入（已废弃）
- **v2** - Garage v2.0.0 引入（当前版本）

## 相关链接

- [Garage 官方文档](https://garagehq.deuxfleurs.fr/documentation/)
- [OpenAPI 规范 (HTML)](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html)
- [OpenAPI 规范 (JSON)](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.json)
- [Garage 源代码](https://git.deuxfleurs.fr/Deuxfleurs/garage)
