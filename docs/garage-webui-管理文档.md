# Garage Web UI Project Management Documentation

## Project Overview

**Garage Web UI** is a modern web management interface for the [Garage](https://garagehq.deuxfleurs.fr/) distributed object storage service. This project provides a clean, intuitive graphical interface to manage Garage clusters, serving as an important supplement to the official Garage command-line tools.

### 🎯 Project Positioning

- **Target Users**: Garage cluster administrators and operations personnel
- **Core Value**: Simplify the daily management operations of Garage clusters
- **Technology Stack**: TypeScript + React (Frontend) + Go (Backend)

## Features

### 🏥 Cluster Monitoring and Management

#### 1. Health Status Monitoring

- **Real-time Cluster Status**: Displays the overall health of the cluster (Healthy/Degraded/Unavailable)
- **Node Monitoring**: Real-time monitoring of the number of known nodes, connected nodes, and storage node status
- **Partition Status**: Monitors the health and quorum status of data partitions

#### 2. Cluster Layout Management

- **Visual Layout**: Graphically displays the cluster node distribution and storage configuration
- **Node Configuration**: Manage node attributes such as zone, capacity, and tags
- **Layout Changes**: Supports staging, previewing, applying, and reverting layout changes
- **History**: View the history of cluster layout changes

### 🗄️ Bucket Management

#### 1. Bucket Operations

- **Bucket List**: Displays all buckets and their basic information
- **Bucket Details**: View detailed statistics, configuration, and permission information for a bucket
- **Bucket Creation**: Supports creating buckets with global and local aliases
- **Bucket Configuration**: Update bucket website configuration, quota settings, etc.

#### 2. Object Browser

- **File Browsing**: Built-in object browser that supports folder structure browsing
- **File Operations**: Upload, download, and delete object files
- **Sharing Functionality**: Generate temporary access links
- **Bulk Operations**: Supports bulk file management

### 🔑 Access Control Management

#### 1. Access Key Management

- **Key List**: Displays all API access keys
- **Key Creation**: Create new S3-compatible access keys
- **Permission Configuration**: Set global permissions for keys (e.g., creating buckets)
- **Expiration Management**: Set expiration times for keys

#### 2. Permission Assignment

- **Bucket Permissions**: Assign permissions to access keys for specific buckets
- **Permission Types**: Supports Read, Write, and Owner permission levels
- **Permission Revocation**: Flexible mechanism for granting and revoking permissions

## Technical Architecture

### 🏗️ Overall Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Browser   │───▶│  Garage Web UI   │───▶│  Garage Cluster │
│  (Frontend UI)  │    │  (Go Backend)    │    │   (Admin API)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 📁 Project Structure

```
garage-webui/
├── src/                    # React Frontend Source
│   ├── pages/             # Page Components
│   │   ├── home/         # Home Dashboard
│   │   ├── cluster/      # Cluster Management
│   │   ├── buckets/      # Bucket Management
│   │   └── keys/         # Access Key Management
│   ├── components/       # Reusable Components
│   ├── hooks/           # React Hooks
│   └── lib/             # Utility Libraries
├── backend/              # Go Backend Source
│   ├── main.go          # Service Entrypoint
│   ├── router/          # API Routes
│   ├── utils/           # Utility Functions
│   └── schema/          # Data Structures
├── docs/                # Project Documentation
└── misc/                # Screenshots and other resources
```

### 🔌 Backend Service Architecture

#### Core Modules

1. **Configuration Management** (`utils/garage.go`)

   - Automatically reads the Garage configuration file (`garage.toml`)
   - Extracts admin API endpoints, authentication information, etc.
   - Supports overriding configuration with environment variables

2. **API Proxy** (`router/`)

   - Proxies frontend requests to the Garage Admin API
   - Handles authentication and error translation
   - Provides a unified RESTful interface

3. **Session Management** (`utils/session.go`)

   - Supports user authentication (optional)
   - Session state management

4. **Caching Mechanism** (`utils/cache.go`)
   - Caches API responses
   - Reduces request pressure on the Garage cluster

## Deployment Scenarios

### 🐳 Docker Deployment (Recommended)

#### 1. Deploying with a Garage Cluster

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

#### 2. Standalone Deployment

```bash
docker run -p 3909:3909 \
  -v ./garage.toml:/etc/garage.toml:ro \
  -e API_BASE_URL="http://garage-host:3903" \
  -e API_ADMIN_KEY="your-admin-token" \
  khairul169/garage-webui:latest
```

### 🖥️ Binary Deployment

```bash
# Download the binary file
wget -O garage-webui https://github.com/khairul169/garage-webui/releases/download/1.0.9/garage-webui-v1.0.9-linux-amd64
chmod +x garage-webui

# Run the service
CONFIG_PATH=./garage.toml ./garage-webui
```

### 🔧 SystemD Service

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

## Configuration Management

### 📝 Garage Configuration Requirements

The Web UI requires the Garage cluster to have the Admin API enabled:

```toml
# garage.toml
[admin]
api_bind_addr = "[::]:3903"
admin_token = "your-secure-admin-token"
metrics_token = "your-metrics-token"  # Optional
```

### 🌍 Environment Variable Configuration

| Variable Name     | Description                   | Default Value        |
| ----------------- | ----------------------------- | -------------------- |
| `CONFIG_PATH`     | Path to Garage config file    | `/etc/garage.toml`   |
| `API_BASE_URL`    | Garage Admin API address      | Read from config file|
| `API_ADMIN_KEY`   | Admin API token               | Read from config file|
| `S3_ENDPOINT_URL` | S3 API address                | Read from config file|
| `S3_REGION`       | S3 region                     | `garage`             |
| `BASE_PATH`       | Web UI base path              | `/`                  |
| `PORT`            | Service port                  | `3909`               |
| `HOST`            | Binding address               | `0.0.0.0`            |

### 🔐 Authentication Configuration

#### Enable Web UI Authentication

```bash
# Generate password hash
htpasswd -nbBC 10 "admin" "password"

# Set environment variable
AUTH_USER_PASS="admin:$2y$10$DSTi9o..."
```

## Management Best Practices

### 🚀 Daily Operations

#### 1. Cluster Health Monitoring

- **Regular Checks**: Monitor cluster status via the home dashboard
- **Alerting Setup**: Configure monitoring systems to connect to the `/metrics` endpoint
- **Performance Observation**: Pay attention to storage node connection status and partition health

#### 2. Bucket Management

- **Naming Conventions**: Establish uniform bucket naming conventions
- **Minimize Permissions**: Assign the minimum necessary permissions to access keys
- **Quota Management**: Set appropriate quota limits for important services

#### 3. Access Control

- **Regular Rotation**: Rotate API access keys periodically
- **Permission Audits**: Regularly review bucket permission assignments
- **Key Management**: Create dedicated access keys for different purposes

### 🔧 Troubleshooting

#### 1. Connection Issues

```bash
# Check Admin API accessibility
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://garage-host:3903/v2/GetClusterHealth

# Check network connectivity
telnet garage-host 3903
```

#### 2. Configuration Issues

- Verify the correctness of the `garage.toml` configuration
- Confirm that the Admin API port is open
- Check firewall and network policies

#### 3. Performance Optimization

- Enable caching to reduce API calls
- Use a reverse proxy (like Nginx) for SSL termination
- Monitor resource usage

### 📊 Monitoring Integration

#### Prometheus Metrics

The Web UI can be configured to monitor Garage's Prometheus metrics:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: "garage"
    static_configs:
      - targets: ["garage-host:3903"]
    metrics_path: /metrics
    bearer_token: "your-metrics-token"
```

#### Key Metrics

- `garage_cluster_health`: Cluster health status
- `garage_storage_usage`: Storage usage
- `garage_api_requests`: API request statistics
- `garage_replication_status`: Data replication status

## Development Guide

### 🛠️ Development Environment Setup

```bash
# Clone the project
git clone https://github.com/khairul169/garage-webui.git
cd garage-webui

# Install frontend dependencies
pnpm install

# Install backend dependencies
cd backend && go mod download && cd ..

# Start the development server
pnpm run dev
```

### 🔧 Technology Choices

- **Frontend**: React 18 + TypeScript + Tailwind CSS
- **State Management**: React Query (TanStack Query)
- **Routing**: React Router
- **UI Components**: Custom component library
- **Backend**: Go + Gin framework
- **Configuration Parsing**: go-toml

### 📋 Contribution Guidelines

1. **Coding Standards**: Follow the project's ESLint and Go fmt standards
2. **Testing**: New features require corresponding tests
3. **Documentation**: Update relevant documents and API descriptions
4. **Compatibility**: Ensure compatibility with the latest version of Garage

## Security Considerations

### 🔒 Security Recommendations

1. **Network Security**

   - Use HTTPS in production environments
   - Restrict network access to the Admin API
   - Use firewall rules to protect sensitive ports

2. **Authentication Security**

   - Consider integrating with an enterprise identity authentication system

3. **Permission Control**
   - Follow the principle of least privilege
   - Use a dedicated administrator token

## Future Plans

### 🚀 Feature Roadmap

- **Advanced Monitoring**: Integrate more performance metrics and alerting features
- **Bulk Operations**: Support bulk management of buckets and access keys
- **API Expansion**: Support more Garage Admin API features
- **Internationalization**: Multi-language support
- **Theming System**: Customizable UI themes

### 🔧 Technical Improvements

- **Cache Optimization**: Smarter caching strategies
- **Real-time Updates**: WebSocket support for real-time status updates
- **Mobile Optimization**: Improve the mobile experience
- **Performance Enhancements**: Frontend bundle optimization and lazy loading

## Garage Admin API Usage

### 🔌 API Features Currently Used by the Project

Based on code analysis, the current Garage Web UI project calls the following Garage Admin API v1 features:

#### 1. Cluster Management API

- **`GET /v1/health`** - Get cluster health status

  - Used on the home dashboard to display cluster status
  - Monitors the number of connected nodes, storage node status, and partition health

- **`GET /v1/status`** - Get detailed cluster status
  - Used on the cluster management page to display node details
  - Shows cluster topology and node configuration information

#### 2. Cluster Layout Management API

- **`GET /v1/layout`** - Get cluster layout configuration

  - Displays the current cluster layout and staged changes
  - Views node roles, capacity, and zone assignments

- **`POST /v1/layout`** - Update cluster layout

  - Adds new nodes to the cluster
  - Modifies node configuration (capacity, zone, tags)
  - Removes nodes (by setting `remove: true`)

- **`POST /v1/connect`** - Connect cluster nodes

  - Connects new nodes to the cluster
  - Establishes RPC connections between nodes

- **`POST /v1/layout/apply`** - Apply layout changes

  - Applies staged layout changes to the cluster
  - Triggers data redistribution

- **`POST /v1/layout/revert`** - Revert layout changes
  - Clears staged layout changes
  - Restores to the last stable state

#### 3. Bucket Management API

- **`GET /v1/bucket?list`** - List all buckets

  - Gets a list of all buckets in the cluster
  - Displays basic bucket information and aliases

- **`GET /v1/bucket?id={id}`** - Get detailed bucket information

  - Views the complete configuration of a single bucket
  - Includes permissions, statistics, quota information, etc.

- **`POST /v1/bucket`** - Create a new bucket

  - Supports setting global and local aliases
  - Configures initial permissions and parameters

- **`PUT /v1/bucket?id={id}`** - Update bucket configuration

  - Modifies the bucket's website configuration
  - Sets or updates quota limits

- **`DELETE /v1/bucket?id={id}`** - Delete a bucket
  - Deletes an empty bucket (the bucket must be empty)

#### 4. Bucket Alias Management API

- **`PUT /v1/bucket/alias/global`** - Add a global alias

  - Creates a global access alias for a bucket
  - Supports multiple aliases pointing to the same bucket

- **`DELETE /v1/bucket/alias/global`** - Delete a global alias
  - Removes a global alias from a bucket
  - The bucket itself remains unaffected

#### 5. Permission Management API

- **`POST /v1/bucket/allow`** - Grant bucket permissions

  - Assigns bucket operation permissions to an access key
  - Supports Read, Write, and Owner permissions

- **`POST /v1/bucket/deny`** - Revoke bucket permissions
  - Removes an access key's permissions for a bucket
  - Flexible permission control mechanism

#### 6. Access Key Management API

- **`GET /v1/key?list`** - List all access keys

  - Gets all API keys in the cluster
  - Displays basic key information

- **`POST /v1/key`** - Create a new access key

  - Generates a new S3-compatible access key
  - Sets initial permissions for the key

- **`POST /v1/key/import`** - Import an existing access key

  - Used for migrating or restoring access keys
  - Imports externally generated keys

- **`DELETE /v1/key?id={id}`** - Delete an access key
  - Removes an access key from the cluster
  - Immediately revokes all related permissions

### ## API Version Comparison Analysis

### 📊 API Differences: Current Project vs. Official Documentation

A comparative analysis reveals that the current project uses **Garage Admin API v1**, while the latest official documentation recommends using **API v2**. Below is a detailed comparison of the differences:

#### 🔄 Version Mapping

| Feature Category         | Current Project (v1)            | Official Recommendation (v2)                          | Status      |
| ---------------- | ------------------------ | -------------------------------------- | --------- |
| **Cluster Health Status** | `GET /v1/health`         | `GET /v2/GetClusterHealth`             | ⚠️ Upgrade Needed |
| **Cluster Status**     | `GET /v1/status`         | `GET /v2/GetClusterStatus`             | ⚠️ Upgrade Needed |
| **Cluster Statistics**     | ❌ Not Used                | `GET /v2/GetClusterStatistics`         | 🆕 New Feature |
| **Connect Nodes**     | `POST /v1/connect`       | `POST /v2/ConnectClusterNodes`         | ⚠️ Upgrade Needed |
| **Get Layout**     | `GET /v1/layout`         | `GET /v2/GetClusterLayout`             | ⚠️ Upgrade Needed |
| **Update Layout**     | `POST /v1/layout`        | `POST /v2/UpdateClusterLayout`         | ⚠️ Upgrade Needed |
| **Apply Layout**     | `POST /v1/layout/apply`  | `POST /v2/ApplyClusterLayout`          | ⚠️ Upgrade Needed |
| **Revert Layout**     | `POST /v1/layout/revert` | `POST /v2/RevertClusterLayout`         | ⚠️ Upgrade Needed |
| **Layout History**     | ❌ Not Used                | `GET /v2/GetClusterLayoutHistory`      | 🆕 New Feature |
| **Preview Layout Changes** | ❌ Not Used                | `POST /v2/PreviewClusterLayoutChanges` | 🆕 New Feature |

#### 📦 Bucket Management API Comparison

| Feature           | Current Project (v1)                    | Official Recommendation (v2)                       | Difference Explanation            |
| -------------- | -------------------------------- | ----------------------------------- | ------------------- |
| **List Buckets** | `GET /v1/bucket?list`            | `GET /v2/ListBuckets`               | Parameter format differs        |
| **Get Bucket Info** | `GET /v1/bucket?id={id}`         | `GET /v2/GetBucketInfo`             | Supports more query methods    |
| **Create Bucket** | `POST /v1/bucket`                | `POST /v2/CreateBucket`             | v2 supports more configuration options |
| **Update Bucket** | `PUT /v1/bucket?id={id}`         | `POST /v2/UpdateBucket/{id}`        | HTTP method and path differ |
| **Delete Bucket** | `DELETE /v1/bucket?id={id}`      | `POST /v2/DeleteBucket/{id}`        | HTTP method differs       |
| **Add Alias**   | `PUT /v1/bucket/alias/global`    | `POST /v2/AddBucketAlias`           | Supports local aliases        |
| **Delete Alias**   | `DELETE /v1/bucket/alias/global` | `POST /v2/RemoveBucketAlias`        | Supports local aliases        |
| **Cleanup Uploads**   | ❌ Not Used                        | `POST /v2/CleanupIncompleteUploads` | 🆕 New Feature           |
| **Inspect Object**   | ❌ Not Used                        | `GET /v2/InspectObject`             | 🆕 New Feature           |

#### 🔑 Access Key Management API Comparison

| Feature             | Current Project (v1)            | Official Recommendation (v2)             | Difference Explanation        |
| ---------------- | ------------------------ | ------------------------- | --------------- |
| **List Keys**     | `GET /v1/key?list`       | `GET /v2/ListKeys`        | Parameter format differs    |
| **Get Key Info** | ❌ Not Used                | `GET /v2/GetKeyInfo`      | 🆕 New Feature       |
| **Create Key**     | `POST /v1/key`           | `POST /v2/CreateKey`      | v2 supports more options |
| **Update Key**     | ❌ Not Used                | `POST /v2/UpdateKey/{id}` | 🆕 New Feature       |
| **Delete Key**     | `DELETE /v1/key?id={id}` | `POST /v2/DeleteKey/{id}` | HTTP method differs   |
| **Import Key**     | `POST /v1/key/import`    | `POST /v2/ImportKey`      | Path structure differs    |
| **Grant Permission**     | `POST /v1/bucket/allow`  | `POST /v2/AllowBucketKey` | Path structure differs    |
| **Revoke Permission**     | `POST /v1/bucket/deny`   | `POST /v2/DenyBucketKey`  | Path structure differs    |

### 🚫 v2-Exclusive Features (Not Used in the Current Project)

#### 1. Admin Token Management

- `GET /v2/ListAdminTokens` - List all admin tokens
- `GET /v2/GetAdminTokenInfo` - Get token information
- `GET /v2/GetCurrentAdminTokenInfo` - Get current token information
- `POST /v2/CreateAdminToken` - Create an admin token
- `POST /v2/UpdateAdminToken/{id}` - Update an admin token
- `POST /v2/DeleteAdminToken/{id}` - Delete an admin token

#### 2. Node Management

- `GET /v2/GetNodeInfo/{node}` - Get node information
- `GET /v2/GetNodeStatistics/{node}` - Get node statistics
- `POST /v2/CreateMetadataSnapshot/{node}` - Create a metadata snapshot
- `POST /v2/LaunchRepairOperation/{node}` - Launch a repair operation

#### 3. Worker Process Management

- `POST /v2/ListWorkers/{node}` - List worker processes
- `POST /v2/GetWorkerInfo/{node}` - Get worker process information
- `POST /v2/GetWorkerVariable/{node}` - Get a worker process variable
- `POST /v2/SetWorkerVariable/{node}` - Set a worker process variable

#### 4. Block Management

- `POST /v2/GetBlockInfo/{node}` - Get block information
- `GET /v2/ListBlockErrors/{node}` - List block errors
- `POST /v2/RetryBlockResync/{node}` - Retry a block resync
- `POST /v2/PurgeBlocks/{node}` - Purge blocks

#### 5. Special Endpoints

- `GET /health` - Quick health check (no authentication required)
- `GET /metrics` - Prometheus metrics
- `GET /check` - On-demand TLS check

### ⚡ Upgrade Impact Analysis

#### 🔴 Key Differences

1. **API Path Structure**

   - v1: Uses query parameters (`?id=xxx`)
   - v2: Uses RESTful paths (`/{id}`)

2. **HTTP Methods**

   - v1: Uses a mix of GET/POST/PUT/DELETE
   - v2: Primarily uses GET/POST

3. **Request/Response Format**

   - v2 provides a more structured data format
   - More detailed error messages and status codes

4. **Feature Completeness**
   - v2 offers more advanced management features
   - Better monitoring and maintenance capabilities

#### 🟡 Compatibility Considerations

- **Backward Compatibility**: The v1 API is still available in the current version (marked as deprecated)
- **Migration Recommendation**: Gradually migrate to the v2 API
- **Feature Enhancement**: Utilize new v2 features to improve user experience

### 📋 Upgrade Recommendations

#### 🎯 Short-Term Plan (1-2 months)

1. **API Version Upgrade**

   - Upgrade core API calls from v1 to v2
   - Update the frontend API client
   - Test for compatibility and functional consistency

2. **Basic Feature Enhancements**
   - Add cluster statistics functionality
   - Implement layout history viewing
   - Support layout change previews

#### 🚀 Medium-Term Plan (3-6 months)

1. **New Feature Integration**

   - Admin token management interface
   - Detailed node information and statistics
   - Object inspection and analysis functionality

2. **Monitoring Enhancements**
   - Integrate Prometheus metrics display
   - Real-time health status monitoring
   - Error and alerting system

#### 🎨 Long-Term Plan (6+ months)

1. **Advanced Management Features**

   - Block management and repair tools
   - Worker process monitoring
   - Automated maintenance tasks

2. **User Experience Optimization**
   - Bulk operations support
   - Real-time data updates
   - Improved mobile adaptation

### 📊 Feature Coverage Analysis

| Feature Category       | v1 Available Features | v2 Total Features | Currently Used | Coverage |
| -------------- | ----------- | --------- | -------- | ------ |
| **Cluster Management**   | 4           | 6         | 2        | 33%    |
| **Layout Management**   | 5           | 7         | 5        | 71%    |
| **Bucket Management** | 7           | 9         | 5        | 56%    |
| **Permission Management**   | 2           | 2         | 2        | 100%   |
| **Key Management**   | 4           | 6         | 4        | 67%    |
| **Advanced Features**   | 0           | 25+       | 0        | 0%     |
| **Overall**       | 22          | 55+       | 18       | 33%    |

**Conclusion**: The current project uses only about 33% of the Garage Admin API's features, leaving significant room for functional expansion.
