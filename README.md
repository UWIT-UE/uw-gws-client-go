# UW Groups Web Service (GWS) Client Library for Go

A Go client library for the University of Washington Groups Web Service API, providing programmatic access to group management, membership operations, and search functionality.

## Official API Documentation

For complete API reference and documentation, see: [UW Groups Web Service API Documentation](https://iam-tools.u.washington.edu/apis/gws/)

## Features

- **Group Management**: Create, read, update, and delete groups
- **Membership Operations**: Add/remove members, check membership status, count members
- **History Tracking**: Retrieve group change history with filtering options
- **Search**: Find groups by name, stem, member, owner, and other criteria
- **Entity Management**: Manage group permissions (admins, updaters, readers, etc.)
- **Effective Membership**: Support for both direct and inherited group membership
- **TLS Authentication**: Client certificate authentication support
- **CLI Tool**: Command-line interface (`gwstool`) for interactive operations

## Installation

```bash
go get github.com/uwit-ue/uw-gws-client-go
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/uwit-ue/uw-gws-client-go/gws"
)

func main() {
    // Create client configuration
    config := gws.DefaultConfig()
    config.APIUrl = "https://groups.uw.edu/group_sws/v3"
    config.CAFile = "/path/to/ca.pem"
    config.ClientCert = "/path/to/client.pem"
    config.ClientKey = "/path/to/client.key"

    // Create client
    client, err := gws.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // Get a group
    group, err := client.GetGroup("u_devtools_admin")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Group: %s (%s)", group.DisplayName, group.ID)
}
```

## Configuration

### Client Configuration

```go
config := &gws.Config{
    APIUrl:        "https://groups.uw.edu/group_sws/v3",  // Production URL
    // APIUrl:     "https://eval.groups.uw.edu/group_sws/v3", // Evaluation URL
    Timeout:       30 * time.Second,
    CAFile:        "/path/to/ca.pem",           // CA certificate
    ClientCert:    "/path/to/client.pem",       // Client certificate
    ClientKey:     "/path/to/client.key",       // Client private key
    Synchronized:  false,                       // Wait for cache updates
    SkipTLSVerify: false,                      // Skip TLS verification (not recommended)
}

client, err := gws.NewClient(config)
```

### Default Configuration

```go
// Start with defaults and customize
config := gws.DefaultConfig()
config.CAFile = "/path/to/ca.pem"
config.ClientCert = "/path/to/client.pem"
config.ClientKey = "/path/to/client.key"
```

## Group Operations

### Get Group Information

```go
group, err := client.GetGroup("u_my_group")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Group: %s\n", group.DisplayName)
fmt.Printf("Description: %s\n", group.Description)
fmt.Printf("Contact: %s\n", group.Contact)
fmt.Printf("Created: %d\n", group.Created)
```

### Create a New Group

```go
newGroup := &gws.Group{
    ID:          "u_my_new_group",
    DisplayName: "My New Group",
    Description: "A group for demonstration purposes",
    Contact:     "admin",
}

// Add administrators (required)
newGroup.Admins.AppendEntityByID("admin1", "admin2")

// Add other permissions (optional)
newGroup.Updaters.AppendEntityByID("user1", "user2")
newGroup.Readers.AppendEntityByID("u_staff_group")

createdGroup, err := client.CreateGroup(newGroup)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created group: %s\n", createdGroup.ID)
```

### Update a Group

```go
// Get existing group
group, err := client.GetGroup("u_my_group")
if err != nil {
    log.Fatal(err)
}

// Modify properties
group.DisplayName = "Updated Group Name"
group.Description = "Updated description"

// Add new admin
group.Admins.AppendEntityByID("newadmin")

// Update the group
updatedGroup, err := client.UpdateGroup(group)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated group: %s\n", updatedGroup.DisplayName)
```

### Delete a Group

```go
err := client.DeleteGroup("u_my_group")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Group deleted successfully")
```

## History Operations

### Get Group History

```go
// Get basic history
history, err := client.GetHistory("u_my_group", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("History entries: %d\n", len(history.Data))
for _, entry := range history.Data {
    fmt.Printf("- %s: %s by %s\n",
        time.Unix(entry.Timestamp/1000, 0).Format("2006-01-02 15:04:05"),
        entry.Description,
        entry.User)
}
```

### Get Filtered History

```go
// Create history options with filters
options := &gws.HistoryOptions{}
options.WithMaxResults(10).
    WithOrder(gws.HistoryOrderDescending).
    WithActivityType(gws.HistoryActivityTypeMembership).
    WithMemberID("user1")

// Get filtered history
history, err := client.GetHistory("u_my_group", options)
if err != nil {
    log.Fatal(err)
}

// Process history entries
for _, entry := range history.Data {
    timestamp := time.Unix(entry.Timestamp/1000, 0)
    fmt.Printf("%s - %s: %s\n",
        timestamp.Format("2006-01-02 15:04:05"),
        entry.Activity,
        entry.Description)

    if entry.ActAs != "" {
        fmt.Printf("  (Acting as: %s)\n", entry.ActAs)
    }
}
```

### History Options

```go
// Time-based filtering (last 30 days)
thirtyDaysAgo := time.Now().AddDate(0, 0, -30).UnixMilli()
options := &gws.HistoryOptions{}
options.WithStartTime(thirtyDaysAgo).
    WithMaxResults(100).
    WithOrder(gws.HistoryOrderAscending)

// Activity type filtering
options.WithActivityType(gws.HistoryActivityTypeACL)        // ACL changes only
options.WithActivityType(gws.HistoryActivityTypeMembership) // Membership changes only

// Member-specific history
options.WithMemberID("user1") // Only changes related to user1
```

## Membership Operations

### Get Group Members

```go
// Get direct members
members, err := client.GetMembership("u_my_group")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Direct members: %d\n", len(*members))
for _, member := range *members {
    fmt.Printf("- %s (%s)\n", member.ID, member.Type)
}

// Get effective members (includes inherited)
effectiveMembers, err := client.GetEffectiveMembership("u_my_group")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Effective members: %d\n", len(*effectiveMembers))
```

### Add Members

```go
// Add individual members
notFound, err := client.AddMembers("u_my_group", "user1", "user2", "user3")
if err != nil {
    log.Fatal(err)
}

if len(notFound) > 0 {
    fmt.Printf("Members not found: %v\n", notFound)
}
```

### Remove Members

```go
// Remove specific members
err := client.DeleteMembers("u_my_group", "user1", "user2")
if err != nil {
    log.Fatal(err)
}

// Remove all members
err = client.DeleteAllMembers("u_my_group")
if err != nil {
    log.Fatal(err)
}
```

### Check Membership

```go
// Check direct membership
isMember, err := client.IsMember("u_my_group", "user1")
if err != nil {
    log.Fatal(err)
}

// Check effective membership
isEffectiveMember, err := client.IsEffectiveMember("u_my_group", "user1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Direct member: %t, Effective member: %t\n", isMember, isEffectiveMember)
```

### Member Counts

```go
// Count direct members
directCount, err := client.MemberCount("u_my_group")
if err != nil {
    log.Fatal(err)
}

// Count effective members
effectiveCount, err := client.EffectiveMemberCount("u_my_group")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Direct: %d, Effective: %d\n", directCount, effectiveCount)
```

### Set Complete Membership

```go
// Create new member list
memberList := gws.NewMemberList()
memberList.AddUWNetIDMembers("user1", "user2", "user3")
memberList.AddGroupMembers("u_staff_group")
memberList.AddDNSMembers("server.example.com")

// Replace entire membership
err := client.SetMembership("u_my_group", memberList)
if err != nil {
    log.Fatal(err)
}
```

## Search Operations

### Basic Search

```go
// Search by group name
search := gws.NewSearch().WithName("*devtools*")
groups, err := client.DoSearch(search)
if err != nil {
    log.Fatal(err)
}

for _, group := range groups {
    fmt.Printf("Found: %s (%s)\n", group.DisplayName, group.ID)
}
```

### Advanced Search

```go
// Search for groups where a user is an effective member
search := gws.NewSearch().
    WithMember("user1").
    InEffectiveMembers()

groups, err := client.DoSearch(search)
if err != nil {
    log.Fatal(err)
}

// Search by owner
search = gws.NewSearch().WithOwner("admin.domain.edu")
groups, err = client.DoSearch(search)

// Search by stem
search = gws.NewSearch().WithStem("u_dept")
groups, err = client.DoSearch(search)

// Search by affiliate
search = gws.NewSearch().WithAffiliate("student")
groups, err = client.DoSearch(search)
```

## Working with Entities

Entities represent different types of identities that can have permissions on groups:

```go
// Add different types of entities to admin list
group.Admins.AppendEntityByID(
    "user1",                    // UWNetID
    "u_admin_group",           // Group
    "server.example.com",      // DNS name
    "user@external.edu",       // EPPN
)

// Work with entity lists
adminIDs := group.Admins.ToIDs()
fmt.Printf("Admin IDs: %v\n", adminIDs)

// Check if entity exists
if group.Admins.Contains("user1") {
    fmt.Println("user1 is an admin")
}

// Remove entity
group.Admins.RemoveEntityByID("old_admin")
```

## Member List Operations

```go
// Create and manipulate member lists
members := gws.NewMemberList()

// Add different member types
members.AddUWNetIDMembers("user1", "user2")
members.AddGroupMembers("u_staff", "u_students")
members.AddDNSMembers("server1.uw.edu", "server2.uw.edu")
members.AddEPPNMembers("external@other.edu")

// Filter by member type
uwnetidMembers := members.Match(gws.MemberTypeUWNetID)
fmt.Printf("UWNetID members: %v\n", uwnetidMembers.ToIDs())

// Convert to comma-separated string
memberString := members.ToCommaString()
fmt.Printf("All members: %s\n", memberString)

// Check membership
if members.Contains("user1") {
    fmt.Println("user1 is in the list")
}
```

## Error Handling

```go
group, err := client.GetGroup("nonexistent_group")
if err != nil {
    if gwsErr, ok := err.(*gws.ErrorResponse); ok {
        fmt.Printf("GWS Error %d: %s\n", gwsErr.Code, gwsErr.Message)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Best Practices

### 1. Connection Management
```go
// Create client once and reuse
client, err := gws.NewClient(config)
if err != nil {
    log.Fatal(err)
}
// Use client for multiple operations
```

### 2. Error Handling
```go
// Always check for errors
if err != nil {
    // Handle appropriately for your application
    log.Printf("GWS operation failed: %v", err)
    return err
}
```

### 3. Synchronized Operations
```go
// Enable synchronized mode for consistency
client.EnableSynchronized()

// Perform write operations
err := client.AddMembers("u_my_group", "user1")

// Disable when not needed (default)
client.DisableSynchronized()
```

**Understanding Synchronized Mode:**

The GWS API uses caching for read operations. When you make changes (create/update/delete groups or modify membership), these changes are processed immediately but may not be visible in subsequent read operations until the cache is updated.

- **Synchronized Mode (slower, consistent)**: Write operations wait for cache propagation before returning. Subsequent reads will immediately see the changes.
- **Default Mode (faster, eventually consistent)**: Write operations return immediately after processing. Changes may not be visible in reads for a short period.

```go
// Example: Ensuring immediate consistency
client.EnableSynchronized()

// Add a member
err := client.AddMembers("u_my_group", "newuser")
if err != nil {
    log.Fatal(err)
}

// This read will definitely see the new member
members, err := client.GetMembership("u_my_group")
// newuser will be in the members list

client.DisableSynchronized() // Return to default for better performance
```

### 4. Batch Operations
```go
// Add multiple members at once instead of individual calls
notFound, err := client.AddMembers("u_my_group", "user1", "user2", "user3")

// Use SetMembership for large membership changes
memberList := gws.NewMemberList()
memberList.AddUWNetIDMembers("user1", "user2", "user3")
err = client.SetMembership("u_my_group", memberList)
```

## gwstool CLI

The repository includes `gwstool`, a command-line interface for GWS operations:

```bash
# Build the tool
go build -o gwstool ./cmd/gwstool

# Configure credentials
gwstool config init --interactive

# Group operations
gwstool group get u_my_group
gwstool group create u_new_group --display-name "New Group" --admin "admin1"
gwstool group history u_my_group --size 10 --order d

# Member operations
gwstool member list u_my_group
gwstool member add u_my_group user1 user2
gwstool member remove u_my_group user1

# Search operations
gwstool search --name "*devtools*"
gwstool search --member "user1" --effective
```

See `cmd/gwstool/README.md` for complete CLI documentation.

## Development Status

### Implemented Features
- ✅ Group CRUD operations
- ✅ Membership management (direct and effective)
- ✅ History retrieval with filtering options
- ✅ Search functionality
- ✅ Entity management
- ✅ TLS client authentication
- ✅ CLI tool (gwstool)

### Planned Features
- 🔄 Group move operations
- 🔄 Affiliate management (get/put/delete)
- 🔄 Additional search filters

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

This project is not officially supported by the University of Washington, but contributions are welcomed via GitHub. Feel free to:

- Report issues or bugs
- Submit feature requests
- Contribute code improvements via pull requests
- Improve documentation

For issues or feature requests, please use the GitHub issue tracker. When contributing code, please ensure your changes include appropriate tests and documentation updates.