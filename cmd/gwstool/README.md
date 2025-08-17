# gwstool

CLI tool for the University of Washington Groups Web Service.

## Installation

```bash
go build -o gwstool ./cmd/gwstool
```

## Configuration

gwstool requires a configuration file with GWS credentials. It supports multiple configuration sources with the following priority order:

1. **Command line config** (highest priority): `--config /path/to/config`
2. **User config**: `~/.config/gwstool/config`
3. **System config**: `/etc/gwstool.conf`

**Note**: At least one configuration file must exist with valid credentials for the tool to function.

Create a config file with your GWS credentials:

```
api_url=https://groups.uw.edu/group_sws/v3
ca_file=/path/to/ca.cert
client_cert=/path/to/client.cert
client_key=/path/to/client.key
timeout=30
```

### Configuration Management

```bash
# Show current configuration and sources
gwstool config show

# Initialize user configuration file
gwstool config init

# Initialize with interactive prompts
gwstool config init --interactive

# Validate configuration
gwstool config validate

# Use a specific config file
gwstool --config /path/to/config.conf group get mygroup
```

## Usage

### Group Operations

```bash
# Get group information
gwstool group get <group-id>

# Create a new group (at least one admin is required)
gwstool group create <group-id> --display-name "My Group" --description "Group description" --admin "erich" --admin "admin2"

# Create with multiple permissions
gwstool group create <group-id> --display-name "My Group" --admin "erich" --updater "user1,user2" --reader "group1"

# Update a group
gwstool group update <group-id> --display-name "New Name"

# Delete a group
gwstool group delete <group-id> --confirm
```

### Member Operations

```bash
# List members of a group
gwstool member list <group-id>

# List effective members (includes inherited)
gwstool member list <group-id> --effective

# Get specific member information
gwstool member get <group-id> <member-id>

# Check if someone is a member
gwstool member check <group-id> <member-id>

# Count members
gwstool member count <group-id>

# Add members
gwstool member add <group-id> <member1> <member2> ...

# Remove members
gwstool member remove <group-id> <member1> <member2> ...

# Clear all members
gwstool member clear <group-id> --confirm
```

### Search Operations

```bash
# Search for groups by name (at least one search parameter required)
gwstool search "mygroup"

# Search with specific parameters
gwstool search --name "mygroup*" --stem "uw_it"

# Search for groups containing a specific member
gwstool search --member "johndoe"

# Search by stem
gwstool search --stem "uw_it"

# Interactive search (prompts for search terms)
# Will prompt for: name, stem, member, owner, instructor, affiliate, scope
# You can skip any field by pressing Enter, but at least one is required
gwstool search --interactive
```
```

### Output Formats

By default, output is in plain text format suitable for bash scripting. Use `--output json` for JSON output:

```bash
# Plain text (default) - one result per line
gwstool member list mygroup

# JSON output
gwstool member list mygroup --output json
```

### Interactive Mode

Use the `--interactive` flag to enable interactive prompts. Interactive mode is most useful for:

- **Configuration setup**: `config init --interactive` (prompts for all config values)
- **Group operations**: `group create --interactive` (prompts for required fields)
- **Destructive operations**: `group delete --interactive`, `member clear --interactive` (confirmation prompts)

```bash
# Interactive configuration setup (recommended for first-time users)
gwstool config init --interactive

# Interactive group creation with prompts for all fields
gwstool group create mygroup --interactive

# Interactive confirmation for destructive operations
gwstool group delete mygroup --interactive
gwstool member clear mygroup --interactive

# Search also supports interactive mode (prompts for all search criteria)
# Will ask for: name, stem, member, owner, instructor, affiliate, scope
gwstool search --interactive
```

## Examples

```bash
# List all members of a group, one per line (good for bash scripting)
gwstool member list uw_staff

# Check if a user is in a group
if gwstool member check uw_staff johndoe >/dev/null 2>&1; then
    echo "johndoe is in uw_staff"
fi

# Count members and store in variable
member_count=$(gwstool member count uw_staff --output json | jq '.count')

# Search for groups and process results
gwstool search --stem "uw_it" | while read group; do
    echo "Processing group: $group"
done
```
