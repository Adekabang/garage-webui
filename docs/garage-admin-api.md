# Garage Admin API Documentation

## Overview

The Garage Administration API is a REST API for programmatically managing a Garage cluster, providing complete functionality for cluster management, bucket management, access control, and more. The current version is v2, and the base API address is typically `http://localhost:3903`.

**⚠️ Important Note**: This documentation reflects the implementation used by the Garage Web UI project. For the most accurate and up-to-date API specifications, please refer to the official documentation at:
- **HTML Documentation**: [https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html)
- **JSON Specification**: [https://garagehq.deuxfleurs.fr/api/garage-admin-v2.json](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.json)

**Current Implementation Status**: The Garage Web UI project has been successfully upgraded to use API v2, implementing 18 core endpoints across cluster management, bucket operations, and access control features. This represents approximately 33% of the total v2 API surface, with significant opportunities for enhanced functionality through additional endpoint implementation.

**Implementation Details**: 
- The project uses both standard v2 endpoints and custom backend endpoints for enhanced functionality
- HTTP methods have been adapted for optimal REST API practices (DELETE for deletions, PUT for additions)
- 9 additional custom endpoints provide object browsing, authentication, and enhanced bucket management
- Total of 27 API endpoints (18 standard v2 + 9 custom) are currently implemented
- **Note**: Some HTTP methods may differ from the official specification due to implementation choices

## Authentication

### Bearer Token Authentication

All API requests must include authentication information in the HTTP header:

```http
Authorization: Bearer <token>
```

### Token Types

1. **User-defined Token** (Recommended)

   - Can be dynamically created and managed
   - Supports scope limitations
   - Supports setting an expiration time
   - Created using the `garage admin-token` command

2. **Master Token** (Deprecated)
   - Specified in the configuration file
   - `admin_token`: For admin endpoint access
   - `metrics_token`: For metrics endpoint access

### Example of Creating a User-defined Token

```bash
garage admin-token create --expires-in 30d \
    --scope ListBuckets,GetBucketInfo,ListKeys,GetKeyInfo,CreateBucket,CreateKey,AllowBucketKey,DenyBucketKey \
    my-token
```

## API Endpoint Categories

### 1. Cluster Management

#### Get Cluster Health

- **Endpoint**: `GET /v2/GetClusterHealth`
- **Description**: Returns the global status of the cluster, including the number of connected nodes, healthy storage nodes, partition status, etc.
- **Response Example**:

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

#### Get Cluster Status

- **Endpoint**: `GET /v2/GetClusterStatus`
- **Description**: Returns detailed cluster status information, including node information and layout configuration.

#### Get Cluster Statistics

- **Endpoint**: `GET /v2/GetClusterStatistics`
- **Description**: Gets cluster-level statistics.

#### Connect Cluster Nodes

- **Endpoint**: `POST /v2/ConnectClusterNodes`
- **Description**: Instructs the current node to connect to other Garage nodes.
- **Request Body**: An array of node addresses `["<node_id>@<net_address>"]`

### 2. Cluster Layout Management

#### Get Cluster Layout

- **Endpoint**: `GET /v2/GetClusterLayout`
- **Description**: Returns the current cluster layout configuration and pending changes.

#### Update Cluster Layout

- **Endpoint**: `POST /v2/UpdateClusterLayout`
- **Description**: Submits cluster layout changes to the staging area.
- **Request Body Example**:

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

#### Apply Layout Changes

- **Endpoint**: `POST /v2/ApplyClusterLayout`
- **Description**: Applies staged layout changes to the cluster.
- **Request Body**: `{"version": <layout_version>}`

#### Preview Layout Changes

- **Endpoint**: `POST /v2/PreviewClusterLayoutChanges`
- **Description**: Previews the impact of layout changes without actually applying them.

#### Revert Layout Changes

- **Endpoint**: `POST /v2/RevertClusterLayout`
- **Description**: Clears all staged layout changes.

#### Get Layout History

- **Endpoint**: `GET /v2/GetClusterLayoutHistory`
- **Description**: Gets the history of cluster layout versions.

### 3. Bucket Management

#### List All Buckets

- **Endpoint**: `GET /v2/ListBuckets`
- **Description**: Returns all buckets and their aliases in the cluster.

#### Get Bucket Information

- **Endpoint**: `GET /v2/GetBucketInfo`
- **Parameters**:
  - `id`: Bucket ID
  - `globalAlias`: Global alias
  - `search`: Search pattern
- **Description**: Gets detailed bucket information, including permissions, statistics, quotas, etc.

#### Create Bucket

- **Endpoint**: `POST /v2/CreateBucket`
- **Request Body Example**:

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

#### Update Bucket

- **Endpoint**: `POST /v2/UpdateBucket/{id}`
- **Description**: Updates a bucket's website configuration and quota settings.
- **Request Body Example**:

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

#### Delete Bucket

- **Endpoint**: `POST /v2/DeleteBucket/{id}` (Official Specification)  
- **Implementation**: `DELETE /v2/DeleteBucket?id={id}` (Garage Web UI)
- **Description**: Deletes an empty bucket (this will delete all associated aliases).
- **Note**: The Garage Web UI uses DELETE method with query parameters for REST compliance

#### Cleanup Incomplete Uploads

- **Endpoint**: `POST /v2/CleanupIncompleteUploads`
- **Request Body**: `{"bucketId": "bucket-id", "olderThanSecs": 86400}`

#### Inspect Object

- **Endpoint**: `GET /v2/InspectObject`
- **Parameters**: `bucketId`, `key`
- **Description**: Gets detailed internal status information for an object.

### 4. Bucket Alias Management

#### Add Bucket Alias

- **Endpoint**: `POST /v2/AddBucketAlias`
- **Description**: Adds a global or local alias for a bucket.

#### Remove Bucket Alias

- **Endpoint**: `POST /v2/RemoveBucketAlias`
- **Description**: Removes a bucket's alias.

### 5. Access Key Management

#### List Access Keys

- **Endpoint**: `GET /v2/ListKeys`
- **Description**: Returns all API access keys.

#### Get Key Information

- **Endpoint**: `GET /v2/GetKeyInfo`
- **Parameters**:
  - `id`: Key ID
  - `search`: Search pattern
  - `showSecretKey`: Whether to return the secret key (default is false).

#### Create Access Key

- **Endpoint**: `POST /v2/CreateKey`
- **Request Body Example**:

```json
{
  "name": "my-key",
  "allow": {
    "createBucket": true
  }
}
```

#### Update Access Key

- **Endpoint**: `POST /v2/UpdateKey/{id}`
- **Description**: Updates a key's name, permissions, and expiration time.

#### Delete Access Key

- **Endpoint**: `POST /v2/DeleteKey/{id}` (Official Specification)
- **Implementation**: `DELETE /v2/DeleteKey?id={id}` (Garage Web UI)
- **Description**: Deletes an access key from the cluster.
- **Note**: The Garage Web UI uses DELETE method with query parameters for REST compliance

#### Import Access Key

- **Endpoint**: `POST /v2/ImportKey`
- **Description**: Imports an existing access key (only for migration and backup recovery).

### 6. Permission Management

#### Grant Permission

- **Endpoint**: `POST /v2/AllowBucketKey`
- **Description**: Grants a key permission to perform operations on a bucket.
- **Request Body Example**:

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

#### Deny Permission

- **Endpoint**: `POST /v2/DenyBucketKey`
- **Description**: Removes a key's permission to perform operations on a bucket.

### 7. Admin API Token Management

#### List Admin Tokens

- **Endpoint**: `GET /v2/ListAdminTokens`

#### Get Token Information

- **Endpoint**: `GET /v2/GetAdminTokenInfo`
- **Parameters**: `id` or `search`

#### Get Current Token Information

- **Endpoint**: `GET /v2/GetCurrentAdminTokenInfo`

#### Create Admin Token

- **Endpoint**: `POST /v2/CreateAdminToken`
- **Request Body Example**:

```json
{
  "name": "my-admin-token",
  "expiration": "2025-12-31T23:59:59Z",
  "scope": ["ListBuckets", "GetBucketInfo", "CreateBucket"]
}
```

#### Update Admin Token

- **Endpoint**: `POST /v2/UpdateAdminToken/{id}`

#### Delete Admin Token

- **Endpoint**: `POST /v2/DeleteAdminToken/{id}`

### 8. Node Management

#### Get Node Information

- **Endpoint**: `GET /v2/GetNodeInfo/{node}`
- **Parameters**: `node` - Node ID, `*` (all nodes), or `self` (current node)

#### Get Node Statistics

- **Endpoint**: `GET /v2/GetNodeStatistics/{node}`

#### Create Metadata Snapshot

- **Endpoint**: `POST /v2/CreateMetadataSnapshot/{node}`

#### Launch Repair Operation

- **Endpoint**: `POST /v2/LaunchRepairOperation/{node}`
- **Repair Types**: `tables`, `blocks`, `versions`, `multipartUploads`, `blockRefs`, `blockRc`, `rebalance`, `aliases`

### 9. Worker Process Management

#### List Workers

- **Endpoint**: `POST /v2/ListWorkers/{node}`

#### Get Worker Information

- **Endpoint**: `POST /v2/GetWorkerInfo/{node}`

#### Get Worker Variable

- **Endpoint**: `POST /v2/GetWorkerVariable/{node}`

#### Set Worker Variable

- **Endpoint**: `POST /v2/SetWorkerVariable/{node}`

### 10. Block Management

#### Get Block Information

- **Endpoint**: `POST /v2/GetBlockInfo/{node}`
- **Request Body**: `{"blockHash": "hash-value"}`

#### List Block Errors

- **Endpoint**: `GET /v2/ListBlockErrors/{node}`

#### Retry Block Resync

- **Endpoint**: `POST /v2/RetryBlockResync/{node}`

#### Purge Blocks

- **Endpoint**: `POST /v2/PurgeBlocks/{node}`
- **Warning**: This operation permanently deletes all objects that reference these blocks.

### 11. Special Endpoints

#### Health Check

- **Endpoint**: `GET /health`
- **Authentication**: None required
- **Description**: Quick health check, returns 200 if the service is available.

#### Prometheus Metrics

- **Endpoint**: `GET /metrics`
- **Authentication**: Optional (using `metrics_token`)
- **Description**: Returns monitoring metrics in Prometheus format.

#### On-Demand TLS Check

- **Endpoint**: `GET /check?domain=<domain>`
- **Authentication**: None required
- **Description**: Used for on-demand TLS certificate validation by reverse proxies (like Caddy).

## Usage Example

### Using curl

```bash
# Get cluster health status
curl -H 'Authorization: Bearer YOUR_TOKEN' \
     http://localhost:3903/v2/GetClusterHealth

# Create a bucket
curl -X POST \
     -H 'Authorization: Bearer YOUR_TOKEN' \
     -H 'Content-Type: application/json' \
     -d '{"globalAlias": "my-bucket"}' \
     http://localhost:3903/v2/CreateBucket

# List all buckets
curl -H 'Authorization: Bearer YOUR_TOKEN' \
     http://localhost:3903/v2/ListBuckets
```

### Using the Garage CLI

```bash
# Call via internal RPC (no authentication required)
garage json-api GetClusterHealth

# Call with parameters
garage json-api GetBucketInfo '{"globalAlias": "my-bucket"}'

# Read parameters from standard input
garage json-api CreateBucket -
{"globalAlias": "test-bucket"}
<EOF>
```

## Error Handling

The API uses standard HTTP status codes:

- `200 OK` - Request successful
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication failed
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Internal server error

Error responses typically include detailed error information:

```json
{
  "error": "Bucket not found",
  "code": "BucketNotFound"
}
```

## Permission Scopes

Admin tokens can be restricted to specific API endpoints:

- `*` - Allows all endpoints
- `ListBuckets` - List buckets
- `GetBucketInfo` - Get bucket information
- `CreateBucket` - Create a bucket
- `ListKeys` - List access keys
- `CreateKey` - Create an access key
- `AllowBucketKey` - Grant permissions
- `DenyBucketKey` - Deny permissions
- `Metrics` - Access the metrics endpoint

## Best Practices

1. **Use User-defined Tokens**: Avoid using the master token from the configuration file.
2. **Set Appropriate Scopes**: Grant only necessary permissions.
3. **Set Expiration Times**: Rotate tokens periodically.
4. **Monitor API Usage**: Monitor API calls via the `/metrics` endpoint.
5. **Handle Errors**: Properly handle various error conditions.
6. **Bulk Operations**: For a large number of operations, consider using bulk APIs or scripts.

## Version History

- **v0** - First introduced in Garage v0.7.2 (deprecated)
- **v1** - Introduced in Garage v0.9.0 (deprecated, still supported for backward compatibility)
- **v2** - Introduced in Garage v2.0.0 (current version, actively used by Garage Web UI)

**Migration Notes**: 
- The Garage Web UI has successfully migrated from v1 to v2 API
- All core functionality now utilizes v2 endpoints for improved performance and feature access
- Legacy v1 endpoints remain available for backward compatibility but are not recommended for new implementations

## Related Links

- [Garage Official Documentation](https://garagehq.deuxfleurs.fr/documentation/)
- [OpenAPI Specification (HTML)](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html)
- [OpenAPI Specification (JSON)](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.json)
- [Garage Source Code](https://git.deuxfleurs.fr/Deuxfleurs/garage)
