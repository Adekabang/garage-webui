# Garage Web UI API Alignment Task List

## Overview

This document outlines the tasks required to fully align the Garage Web UI implementation with the official [Garage Admin API v2 specification](https://garagehq.deuxfleurs.fr/api/garage-admin-v2.html). The current implementation uses REST-compliant HTTP methods but may differ from the official specification.

**Current Status**: 18/55+ official v2 endpoints implemented (33% coverage)  
**Goal**: Achieve full compliance with official API specification while maintaining REST best practices

---

## 🔧 **High Priority: HTTP Method Alignment**

### Task 1: Verify and Align Delete Operations
- [ ] **Research Official Specification**: Confirm the exact HTTP methods specified for delete operations in the official docs
- [ ] **Update DeleteKey Implementation**: 
  - Current: `DELETE /v2/DeleteKey?id={id}`
  - Official (likely): `POST /v2/DeleteKey/{id}`
  - Decision needed: Keep REST-compliant DELETE or align with official POST
- [ ] **Update DeleteBucket Implementation**:
  - Current: `DELETE /v2/DeleteBucket?id={id}`
  - Official (likely): `POST /v2/DeleteBucket/{id}`
  - Decision needed: Keep REST-compliant DELETE or align with official POST
- [ ] **Update Frontend Hooks**: Modify `src/pages/keys/hooks.ts` and `src/pages/buckets/manage/hooks.ts` if changes are made
- [ ] **Test Compatibility**: Ensure changes work with actual Garage server instances

### Task 2: Review Other HTTP Methods
- [ ] **Verify UpdateBucket Method**: Confirm if `POST /v2/UpdateBucket` is correct vs potential `PUT`
- [ ] **Check Alias Operations**: Verify `PUT/DELETE` methods for alias operations are officially correct
- [ ] **Validate All POST Operations**: Ensure all POST endpoints match official specification

---

## 📋 **Medium Priority: Missing Official Endpoints**

### Task 3: Implement Additional Cluster Management
- [ ] **Add GetClusterStatistics**: `GET /v2/GetClusterStatistics`
  - Create hook in `src/pages/cluster/hooks.ts`
  - Add UI component for cluster statistics display
  - Integrate into cluster dashboard
- [ ] **Add PreviewClusterLayoutChanges**: `POST /v2/PreviewClusterLayoutChanges`
  - Implement preview functionality before applying changes
  - Add confirmation dialog with preview results
- [ ] **Add GetClusterLayoutHistory**: `GET /v2/GetClusterLayoutHistory`
  - Create layout history viewer component
  - Add navigation to view past layout versions

### Task 4: Implement Enhanced Bucket Features
- [ ] **Add CleanupIncompleteUploads**: `POST /v2/CleanupIncompleteUploads`
  - Create cleanup functionality in bucket management
  - Add scheduled/manual cleanup options
- [ ] **Add InspectObject**: `GET /v2/InspectObject`
  - Integrate into object browser
  - Add object inspection modal/page
- [ ] **Add Enhanced Bucket Operations**:
  - [ ] `POST /v2/AddBucketAlias` (if different from current implementation)
  - [ ] `POST /v2/RemoveBucketAlias` (if different from current implementation)

### Task 5: Implement Key Management Enhancements
- [ ] **Add GetKeyInfo**: `GET /v2/GetKeyInfo`
  - Create key details page
  - Show key capabilities, expiration, etc.
- [ ] **Add UpdateKey**: `POST /v2/UpdateKey/{id}`
  - Add key editing functionality
  - Allow updating permissions, expiration, name

---

## 🚀 **Low Priority: Advanced Features**

### Task 6: Admin Token Management
- [ ] **Add Token Listing**: `GET /v2/ListAdminTokens`
- [ ] **Add Token Details**: `GET /v2/GetAdminTokenInfo`
- [ ] **Add Current Token Info**: `GET /v2/GetCurrentAdminTokenInfo`
- [ ] **Add Token Creation**: `POST /v2/CreateAdminToken`
- [ ] **Add Token Updates**: `POST /v2/UpdateAdminToken/{id}`
- [ ] **Add Token Deletion**: `POST /v2/DeleteAdminToken/{id}`
- [ ] **Create Admin Token Management UI**:
  - Token listing page
  - Token creation wizard
  - Token permissions management

### Task 7: Node Management & Monitoring
- [ ] **Add Node Information**: `GET /v2/GetNodeInfo/{node}`
- [ ] **Add Node Statistics**: `GET /v2/GetNodeStatistics/{node}`
- [ ] **Add Metadata Snapshots**: `POST /v2/CreateMetadataSnapshot/{node}`
- [ ] **Add Repair Operations**: `POST /v2/LaunchRepairOperation/{node}`
- [ ] **Create Node Management UI**:
  - Node dashboard with detailed information
  - Repair operation scheduler
  - Node health monitoring

### Task 8: Worker Process Management
- [ ] **Add Worker Listing**: `POST /v2/ListWorkers/{node}`
- [ ] **Add Worker Information**: `POST /v2/GetWorkerInfo/{node}`
- [ ] **Add Worker Variables**: `POST /v2/GetWorkerVariable/{node}`
- [ ] **Add Variable Setting**: `POST /v2/SetWorkerVariable/{node}`
- [ ] **Create Worker Management UI**:
  - Worker process monitor
  - Variable configuration interface

### Task 9: Block Management
- [ ] **Add Block Information**: `POST /v2/GetBlockInfo/{node}`
- [ ] **Add Block Error Listing**: `GET /v2/ListBlockErrors/{node}`
- [ ] **Add Block Resync**: `POST /v2/RetryBlockResync/{node}`
- [ ] **Add Block Purging**: `POST /v2/PurgeBlocks/{node}`
- [ ] **Create Block Management UI**:
  - Block health dashboard
  - Error resolution tools
  - Maintenance operations interface

---

## 🧪 **Testing & Validation Tasks**

### Task 10: Compatibility Testing
- [ ] **Test Against Multiple Garage Versions**:
  - [ ] Test with Garage v2.0.x
  - [ ] Test with latest Garage version
  - [ ] Verify backward compatibility
- [ ] **API Method Validation**:
  - [ ] Test all current endpoints with official methods
  - [ ] Verify error handling consistency
  - [ ] Check response format compliance
- [ ] **Integration Testing**:
  - [ ] Test with real Garage clusters
  - [ ] Validate cluster operations end-to-end
  - [ ] Test authentication and authorization

### Task 11: Error Handling Alignment
- [ ] **Review Error Response Formats**: Ensure they match official specification
- [ ] **Update Error Messages**: Align with official API error codes and messages
- [ ] **Implement Proper Status Codes**: Verify all HTTP status codes match specification

---

## 📚 **Documentation & Maintenance Tasks**

### Task 12: Documentation Updates
- [ ] **Update API Documentation**: Align all endpoint documentation with official specification
- [ ] **Add Implementation Notes**: Document any intentional deviations from official spec
- [ ] **Create Migration Guide**: Document changes needed for users upgrading
- [ ] **Update README**: Reflect full API v2 compliance status

### Task 13: Code Quality Improvements
- [ ] **Type Definition Updates**: Create proper TypeScript interfaces for all API responses
- [ ] **Error Handling Standardization**: Implement consistent error handling across all endpoints
- [ ] **Code Documentation**: Add comprehensive JSDoc comments to all API functions
- [ ] **Unit Testing**: Add tests for all API integration functions

---

## 🎯 **Implementation Strategy**

### Phase 1: Core Alignment (Week 1-2)
1. Complete Tasks 1-2 (HTTP Method Alignment)
2. Implement Task 10 (Compatibility Testing)
3. Update documentation for changes

### Phase 2: Essential Features (Week 3-4)
1. Complete Task 3 (Additional Cluster Management)
2. Complete Task 4 (Enhanced Bucket Features)
3. Complete Task 5 (Key Management Enhancements)

### Phase 3: Advanced Features (Month 2)
1. Implement Tasks 6-9 based on user demand and priority
2. Complete comprehensive testing
3. Full documentation update

### Phase 4: Maintenance (Ongoing)
1. Monitor official API updates
2. Maintain compatibility with new Garage versions
3. Address user feedback and issues

---

## ✅ **Completion Criteria**

- [ ] All official v2 endpoints implemented (100% coverage)
- [ ] HTTP methods align with official specification (or documented deviations)
- [ ] Comprehensive test coverage for all endpoints
- [ ] Updated documentation reflecting full compliance
- [ ] Successful testing against multiple Garage versions
- [ ] User feedback incorporated and addressed

---

**Last Updated**: July 2025  
**Next Review**: Monthly or when official API specification changes 
