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
    image: adekabang/garage-webui:latest
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
  adekabang/garage-webui:latest
```

### 🖥️ Binary Deployment

```bash
# Download the binary file
wget -O garage-webui https://github.com/Adekabang/garage-webui/releases/download/1.0.9/garage-webui-v1.0.9-linux-amd64
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
git clone https://github.com/Adekabang/garage-webui.git
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

## API Features Currently Used by the Project

**✅ Full API Compliance**: As of July 2025, all API implementations are fully aligned with the official [Garage Admin API v2 specification](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html). All HTTP methods, request formats, and response handling match the official documentation exactly.

#### 1. Cluster Management API

- **`GET /v2/GetClusterHealth`** - Get cluster health status

  - Used on the home dashboard to display cluster status
  - Monitors the number of connected nodes, storage node status, and partition health

- **`GET /v2/GetClusterStatus`** - Get detailed cluster status
  - Used on the cluster management page to display node details
  - Shows cluster topology and node configuration information

#### 2. Cluster Layout Management API

- **`GET /v2/GetClusterLayout`** - Get cluster layout configuration

  - Displays the current cluster layout and staged changes
  - Views node roles, capacity, and zone assignments

- **`POST /v2/UpdateClusterLayout`** - Update cluster layout

  - Adds new nodes to the cluster
  - Modifies node configuration (capacity, zone, tags)
  - Removes nodes (by setting `remove: true`)

- **`POST /v2/ConnectClusterNodes`** - Connect cluster nodes

  - Connects new nodes to the cluster
  - Establishes RPC connections between nodes

- **`POST /v2/ApplyClusterLayout`** - Apply layout changes

  - Applies staged layout changes to the cluster
  - Triggers data redistribution

- **`POST /v2/RevertClusterLayout`** - Revert layout changes
  - Clears staged layout changes
  - Restores to the last stable state

#### 3. Bucket Management API

- **`GET /v2/ListBuckets`** - List all buckets

  - Gets a list of all buckets in the cluster
  - Displays basic bucket information and aliases

- **`GET /v2/GetBucketInfo`** - Get detailed bucket information

  - Views the complete configuration of a single bucket
  - Includes permissions, statistics, quota information, etc.

- **`POST /v2/CreateBucket`** - Create a new bucket

  - Supports setting global and local aliases
  - Configures initial permissions and parameters

- **`POST /v2/UpdateBucket`** - Update bucket configuration

  - Modifies the bucket's website configuration
  - Sets or updates quota limits

- **`POST /v2/DeleteBucket?id={id}`** - Delete bucket
  - **Implementation**: POST with query parameter (aligned with official specification)
  - **Parameters**: Bucket ID in query parameter
  - **Usage**: Remove buckets from cluster

#### 4. Bucket Alias Management API

- **`POST /v2/AddBucketAlias`** - Add a global alias
  - **Implementation**: POST with JSON body (aligned with official specification)
  - **Parameters**: Bucket ID and alias name in request body (`{ bucketId, globalAlias }`)
  - **Usage**: Add new bucket aliases

- **`POST /v2/RemoveBucketAlias`** - Remove a global alias  
  - **Implementation**: POST with JSON body (aligned with official specification)
  - **Parameters**: Bucket ID and alias name in request body (`{ bucketId, globalAlias }`)
  - **Usage**: Remove existing bucket aliases

#### 5. Permission Management API

- **`POST /v2/AllowBucketKey`** - Grant bucket permissions

  - Assigns bucket operation permissions to an access key
  - Supports Read, Write, and Owner permissions

- **`POST /v2/DenyBucketKey`** - Revoke bucket permissions
  - Removes an access key's permissions for a bucket
  - Flexible permission control mechanism

#### 6. Access Key Management API

- **`GET /v2/ListKeys`** - List all access keys

  - Gets all API keys in the cluster
  - Displays basic key information

- **`POST /v2/CreateKey`** - Create a new access key

  - Generates a new S3-compatible access key
  - Sets initial permissions for the key

- **`POST /v2/ImportKey`** - Import an existing access key

  - Used for migrating or restoring access keys
  - Imports externally generated keys

- **`POST /v2/DeleteKey?id={id}`** - Delete access key  
  - **Implementation**: POST with query parameter (aligned with official specification)
  - **Parameters**: Key ID in query parameter
  - **Usage**: Remove access keys from cluster

### 🔧 Custom Backend Endpoints

The Garage Web UI implements several custom backend endpoints that extend functionality beyond the standard Garage Admin API:

#### 1. Configuration Management

- **`GET /config`** - Get garage configuration
  - Retrieves garage configuration for frontend display
  - Provides S3 endpoint URLs, region information, etc.

#### 2. Authentication System

- **`POST /auth/login`** - User authentication
  - Handles user login with username/password
  - Creates authenticated sessions

- **`GET /auth/status`** - Authentication status
  - Checks current authentication state
  - Returns whether authentication is enabled and user status

- **`POST /auth/logout`** - User logout
  - Terminates authenticated sessions
  - Clears session data

#### 3. Enhanced Bucket Operations

- **`GET /buckets`** - Enhanced bucket listing
  - Provides enriched bucket information by combining `/v2/ListBuckets` and `/v2/GetBucketInfo`
  - Includes detailed statistics and metadata for all buckets

#### 4. Object Browser & File Management

- **`GET /browse/{bucket}`** - Browse bucket objects
  - Lists objects and folders in a bucket with S3 ListObjectsV2
  - Supports pagination and prefix filtering
  - Provides object metadata and download URLs

- **`GET /browse/{bucket}/{key...}`** - Get/view object
  - Retrieves object content for viewing or downloading
  - Supports thumbnail generation for images
  - Provides object metadata via HeadObject

- **`PUT /browse/{bucket}/{key...}`** - Upload object
  - Handles file uploads with multipart form data
  - Supports directory creation
  - Uses S3 PutObject for storage

- **`DELETE /browse/{bucket}/{key...}`** - Delete object/folder
  - Deletes individual objects or entire folders recursively
  - Supports bulk deletion for folders
  - Uses S3 DeleteObject/DeleteObjects

### 📊 Current Feature Coverage Analysis

| Feature Category       | v2 Total Features | Currently Implemented | Custom Extensions | Total Coverage |
| ---------------------- | ----------------- | -------------------- | ---------------- | -------------- |
| **Cluster Management** | 6                 | 2                    | 1 (config)       | 50%            |
| **Layout Management**  | 7                 | 5                    | 0                | 71%            |
| **Bucket Management**  | 9                 | 5                    | 2 (enhanced list, browse) | 78% |
| **Permission Management** | 2              | 2                    | 0                | 100%           |
| **Key Management**     | 6                 | 4                    | 0                | 67%            |
| **Authentication**     | 0                 | 0                    | 3 (login system) | 100% (custom)  |
| **File Management**    | 0                 | 0                    | 4 (object browser) | 100% (custom) |
| **Advanced Features**  | 25+               | 0                    | 0                | 0%             |
| **Overall**            | 55+               | 18                   | 10               | 51%            |

**Enhanced Coverage**: With custom backend endpoints, the project achieves 51% total feature coverage, providing comprehensive cluster management, authentication, and file browsing capabilities.

### 🚀 Available v2 Features Not Yet Implemented

The project can further enhance functionality by implementing these additional v2 API features:

#### 1. Admin Token Management
- `GET /v2/ListAdminTokens` - List all admin tokens
- `GET /v2/GetAdminTokenInfo` - Get token information
- `GET /v2/GetCurrentAdminTokenInfo` - Get current token information
- `POST /v2/CreateAdminToken` - Create an admin token
- `POST /v2/UpdateAdminToken/{id}` - Update an admin token
- `POST /v2/DeleteAdminToken/{id}` - Delete an admin token

#### 2. Enhanced Node Management
- `GET /v2/GetNodeInfo/{node}` - Get node information
- `GET /v2/GetNodeStatistics/{node}` - Get node statistics
- `POST /v2/CreateMetadataSnapshot/{node}` - Create a metadata snapshot
- `POST /v2/LaunchRepairOperation/{node}` - Launch a repair operation

#### 3. Worker Process Management
- `POST /v2/ListWorkers/{node}` - List worker processes
- `POST /v2/GetWorkerInfo/{node}` - Get worker process information
- `POST /v2/GetWorkerVariable/{node}` - Get a worker process variable
- `POST /v2/SetWorkerVariable/{node}` - Set a worker process variable

#### 4. Advanced Block Management
- `POST /v2/GetBlockInfo/{node}` - Get block information
- `GET /v2/ListBlockErrors/{node}` - List block errors
- `POST /v2/RetryBlockResync/{node}` - Retry a block resync
- `POST /v2/PurgeBlocks/{node}` - Purge blocks

#### 5. Enhanced Bucket Features
- `POST /v2/CleanupIncompleteUploads` - Cleanup incomplete uploads
- `GET /v2/InspectObject` - Inspect object details
- `GET /v2/GetClusterStatistics` - Get cluster-level statistics
- `POST /v2/PreviewClusterLayoutChanges` - Preview layout changes
- `GET /v2/GetClusterLayoutHistory` - Get layout history

### 🎯 Future Development Roadmap

#### 📅 Short-Term (1-2 months)
1. **Enhanced Monitoring**: Implement cluster statistics and node information display
2. **Layout Improvements**: Add layout history viewing and change preview functionality
3. **Object Management**: Add object inspection and incomplete upload cleanup

#### 📅 Medium-Term (3-6 months)
1. **Admin Token Management**: Full admin token lifecycle management interface
2. **Advanced Monitoring**: Worker process monitoring and detailed node statistics
3. **Maintenance Tools**: Automated repair operations and block management

#### 📅 Long-Term (6+ months)
1. **Complete API Coverage**: Implement all available v2 API endpoints
2. **Advanced Bulk Operations**: Comprehensive bulk management features
3. **Real-time Integration**: WebSocket support for live updates and monitoring

**Current Status**: The Garage Web UI successfully uses Garage Admin API v2 with 18 standard endpoints plus 10 custom backend endpoints, achieving 51% feature coverage. The project provides a comprehensive management interface with cluster administration, authentication, file browsing, and enhanced bucket management capabilities, serving as a robust web-based alternative to command-line tools.

**Implementation Philosophy**: The project prioritizes REST API compliance and user experience, which may result in HTTP method choices that differ from the official specification while maintaining full functional compatibility.
