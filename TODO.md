# TODO - Incomplete Features and Known Issues

This document tracks incomplete features, known limitations, and planned enhancements for the UW Groups Web Service Go client library.

## High Priority - Core API Features

### 1. Group Move Operations
**Status:** Not implemented
**Location:** `gws/gws_groups.go:385`
**Description:** The ability to move/rename groups to different paths in the group hierarchy.

**Details:**
- Currently missing the API endpoint implementation for group moves
- This is a common administrative operation needed for group management
- Would likely involve a new endpoint like `PUT /group/{old-id}/move` with new path in body

**Estimated Effort:** Medium - requires new API endpoint and proper validation

### 2. Affiliate Management (CRUD Operations)
**Status:** Partially implemented (search only)
**Location:** `gws/gws_groups.go:386`
**Description:** Full CRUD operations for group affiliates are missing.

**Current Status:**
- ✅ Search by affiliate (`WithAffiliate()` in search)
- ✅ Affiliate type definition in `gws_objects.go`
- ❌ Get affiliate information for a group
- ❌ Add/update affiliate associations
- ❌ Delete affiliate associations

**Details:**
- Affiliates represent different identity types: email, google, uwnetid, radius
- Need endpoints for managing affiliate relationships with groups
- Should support operations like: `GetAffiliates()`, `AddAffiliate()`, `UpdateAffiliate()`, `DeleteAffiliate()`

**Estimated Effort:** Medium-High - requires multiple new endpoints and data structures

### 3. Synchronized Mode for Affiliate Operations
**Status:** Not implemented
**Location:** `gws/gws.go:133`
**Description:** The synchronized mode (wait for cache propagation) doesn't support PUT affiliate operations.

**Details:**
- Currently `syncQueryString()` is only used for group and membership operations
- Affiliate PUT operations should also support synchronized mode for consistency
- Need to add synchronized parameter support to future affiliate endpoints

**Estimated Effort:** Low - extend existing synchronized functionality

## Medium Priority - Enhancements

### 4. Additional Search Filters
**Status:** Enhancement opportunity
**Description:** While basic search is implemented, there may be additional search capabilities in the GWS API not yet exposed.

**Current Search Features:**
- ✅ Name/wildcard search
- ✅ Stem search
- ✅ Member search (direct/effective)
- ✅ Owner search
- ✅ Instructor search
- ✅ Affiliate search
- ✅ Scope filtering

**Potential Enhancements:**
- Advanced date-based filtering
- Complex boolean search combinations
- Additional metadata-based searches
- Performance optimizations for large result sets

**Estimated Effort:** Low-Medium - depends on available API capabilities

### 5. Enhanced Error Handling
**Status:** Improvement opportunity
**Description:** While basic error handling exists, it could be enhanced with more specific error types.

**Current State:**
- Basic error response handling in `gws_errors.go`
- Generic error formatting

**Potential Improvements:**
- Specific error types for different API errors (NotFound, Unauthorized, etc.)
- Better error context and suggestions
- Retry mechanisms for transient failures
- Rate limiting handling

**Estimated Effort:** Medium - requires API error analysis and type definitions

## Low Priority - Nice to Have

### 6. Caching and Performance Optimizations
**Status:** Future enhancement
**Description:** Client-side caching and performance improvements.

**Ideas:**
- Optional client-side caching for frequently accessed groups
- Batch operations for multiple group operations
- Connection pooling optimizations
- Compression support

**Estimated Effort:** High - requires careful design to avoid cache consistency issues

### 7. Advanced CLI Features
**Status:** Enhancement opportunity
**Description:** Additional gwstool capabilities for power users.

**Potential Features:**
- Bulk operations from CSV/JSON files
- Interactive group browsing/navigation
- Advanced output formatting options
- Configuration profiles for different environments
- Shell completion support

**Estimated Effort:** Medium - mostly CLI enhancements, no core library changes needed

## Documentation and Testing

### 8. API Coverage Analysis
**Status:** Ongoing need
**Description:** Systematic review of GWS API to ensure all endpoints are covered.

**Action Items:**
- Compare implemented endpoints against full GWS API documentation
- Identify any missing endpoints or parameters
- Prioritize implementation based on common use cases

### 9. Integration Tests
**Status:** Could be improved
**Description:** More comprehensive testing against live/mock GWS API.

**Current State:**
- Basic examples in `gws-test.go`
- Unit tests could be expanded

**Improvements Needed:**
- Mock server for testing
- Comprehensive integration test suite
- Performance testing
- Error condition testing

## Implementation Notes

### Code Organization
- New affiliate operations should go in a new file `gws/gws_affiliate.go`
- Group move operations should be added to `gws/gws_groups.go`
- Follow existing patterns for API client methods and error handling

### API Design Principles
- Maintain consistency with existing function signatures
- Use appropriate Go idioms (pointer receivers, error returns)
- Support both simple and advanced use cases
- Ensure backward compatibility

### Testing Strategy
- All new features should include unit tests
- Integration tests for critical paths
- Documentation examples should be tested
- Consider adding fuzzing for input validation

## Contributing

When implementing features from this TODO list:

1. Check the [official GWS API documentation](https://iam-tools.u.washington.edu/apis/gws/) for endpoint details
2. Follow existing code patterns and conventions
3. Add comprehensive tests and documentation
4. Update this TODO list to reflect completed work
5. Consider backward compatibility impacts

## Last Updated

This TODO list was generated on August 26, 2025, based on analysis of the codebase at that time. As features are implemented or new requirements discovered, this document should be updated accordingly.
