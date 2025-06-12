# Fast Navigation (fn) GO Application Plan

## Overview
A command-line tool that allows users to bookmark directories and quickly navigate to them using short aliases.

**Core Concept**: `fn save foo` saves current directory as "foo", `fn foo` navigates to that directory.

## Technical Architecture

### Key Challenge & Solution
**Challenge**: A GO binary cannot directly change the parent shell's working directory.

**Solution**: The GO binary outputs the target directory path, and a shell function wrapper evaluates it:
```bash
fn() {
    if [[ "$1" == "save" ]] || [[ "$1" == "list" ]] || [[ "$1" == "delete" ]] || [[ "$1" == "path" ]]; then
        command fn "$@"
    else
        local dir=$(command fn navigate "$@")
        if [[ -n "$dir" ]]; then
            cd "$dir"
        fi
    fi
}
```

## Project Structure
```
fn/
├── main.go
├── cmd/
│   ├── root.go      # Root command setup
│   ├── save.go      # Save current directory
│   ├── list.go      # List all bookmarks
│   ├── delete.go    # Delete a bookmark
│   ├── navigate.go  # Navigate to bookmark (outputs path)
│   └── path.go      # Print path without navigating
├── internal/
│   ├── config/
│   │   └── config.go    # Configuration management
│   └── storage/
│       └── storage.go   # Bookmark storage operations
├── scripts/
│   └── install.sh       # Installation script
├── go.mod
├── go.sum
└── README.md
```

## Core Features

### Commands
1. **`fn save <alias>`** - Save current directory with an alias
   - Validates alias name (no spaces, special chars)
   - Overwrites if alias exists (with confirmation)
   - Stores absolute path

2. **`fn <alias>`** - Navigate to saved directory
   - Outputs directory path for shell to cd into
   - Error if alias doesn't exist
   - Validates directory still exists

3. **`fn list`** - List all saved aliases
   - Formatted output with colors
   - Shows alias → path mapping
   - Indicates if path no longer exists

4. **`fn delete <alias>`** - Remove a saved alias
   - Confirmation prompt
   - Error if alias doesn't exist

5. **`fn path <alias>`** - Print path without navigating
   - Useful for scripts or inspection

6. **`fn edit <alias>`** - Update existing alias to current directory

## Data Storage

### Location
- Config directory: `~/.fn/`
- Bookmarks file: `~/.fn/bookmarks.json`

### JSON Structure
```json
{
  "version": "1.0",
  "bookmarks": {
    "foo": {
      "path": "/home/user/projects/foo",
      "created": "2024-01-15T10:30:00Z",
      "used_count": 42,
      "last_used": "2024-01-20T15:45:00Z"
    },
    "docs": {
      "path": "/home/user/Documents",
      "created": "2024-01-10T09:00:00Z",
      "used_count": 15,
      "last_used": "2024-01-19T11:20:00Z"
    }
  }
}
```

## Implementation Plan

### Phase 1: Core Functionality (MVP)
1. **Project Setup**
   - Initialize GO module
   - Add cobra dependency
   - Basic project structure

2. **Storage Layer**
   - Create/read config directory
   - Load/save bookmarks JSON
   - Handle file permissions

3. **Basic Commands**
   - Implement `save` command
   - Implement `navigate` command (outputs path)
   - Basic error handling

4. **Shell Integration**
   - Create bash function
   - Create zsh function
   - Installation instructions

### Phase 2: Enhanced Features
1. **Additional Commands**
   - `list` with formatted output
   - `delete` with confirmation
   - `edit` command
   - `path` command

2. **Validation & Safety**
   - Validate alias names
   - Check directory existence
   - Handle permission errors
   - Backup before overwrites

3. **User Experience**
   - Colored output
   - Progress indicators
   - Better error messages
   - Success confirmations

### Phase 3: Advanced Features
1. **Autocomplete**
   - Bash completion script
   - Zsh completion script
   - Dynamic alias suggestions

2. **Search & Filter**
   - Fuzzy search aliases
   - Filter by path pattern
   - Sort by usage/recency

3. **Import/Export**
   - Export bookmarks
   - Import from other tools
   - Merge capabilities

### Phase 4: Polish & Optimization
1. **Performance**
   - Optimize JSON parsing
   - Cache for autocomplete
   - Lazy loading

2. **Cross-Platform**
   - Windows support
   - Path normalization
   - Shell detection

3. **Documentation**
   - Man page
   - Examples
   - Troubleshooting guide

## Technical Decisions

### Dependencies
```go
// go.mod
module github.com/username/fn

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/fatih/color v1.16.0
    github.com/mitchellh/go-homedir v1.1.0
    github.com/AlecAivazis/survey/v2 v2.3.7  // For interactive prompts
)
```

### Error Handling Strategy
- User-friendly error messages
- Exit codes:
  - 0: Success
  - 1: General error
  - 2: Alias not found
  - 3: Directory not found
  - 4: Permission denied

### Alias Validation Rules
- Alphanumeric + dash/underscore only
- No spaces or special characters
- Case-sensitive
- Max length: 50 characters
- Reserved words: save, list, delete, edit, path, help

## Installation Process

### Manual Installation
```bash
# Build
go build -o fn

# Install binary
sudo cp fn /usr/local/bin/

# Add to ~/.bashrc
echo 'fn() {
    if [[ "$1" == "save" ]] || [[ "$1" == "list" ]] || [[ "$1" == "delete" ]] || [[ "$1" == "path" ]]; then
        command fn "$@"
    else
        local dir=$(command fn navigate "$@")
        if [[ -n "$dir" ]]; then
            cd "$dir"
        fi
    fi
}' >> ~/.bashrc

# Reload shell
source ~/.bashrc
```

### Automated Installation
Create `install.sh` script that:
1. Detects OS and architecture
2. Downloads appropriate binary
3. Installs to PATH
4. Adds shell function
5. Sets up completions

## Usage Examples

```bash
# Save current directory
$ cd /home/user/projects/myapp
$ fn save myapp
✓ Saved 'myapp' → /home/user/projects/myapp

# Navigate to saved directory
$ cd ~
$ fn myapp
$ pwd
/home/user/projects/myapp

# List all bookmarks
$ fn list
📍 myapp     → /home/user/projects/myapp (used 42 times)
📍 docs      → /home/user/Documents (used 15 times)
📍 downloads → /home/user/Downloads (used 8 times)

# Delete a bookmark
$ fn delete downloads
? Are you sure you want to delete 'downloads'? (y/N) y
✓ Deleted 'downloads'

# Get path without navigating
$ fn path myapp
/home/user/projects/myapp

# Edit existing bookmark
$ cd /home/user/projects/myapp-v2
$ fn edit myapp
✓ Updated 'myapp' → /home/user/projects/myapp-v2
```

## Future Enhancements

### Version 2.0
- **Tags/Categories**: Organize bookmarks with tags
- **Fuzzy Finder Integration**: Use fzf for interactive selection
- **Directory History**: Track navigation patterns
- **Frecency**: Sort by frequency + recency

### Version 3.0
- **Team Sync**: Share bookmarks across team
- **Cloud Backup**: Optional cloud sync
- **Plugins**: Extensible architecture
- **GUI**: Optional system tray app

## Development Guidelines

### Code Style
- Follow standard GO formatting
- Meaningful variable names
- Comments for exported functions
- Unit tests for core logic

### Testing Strategy
- Unit tests for storage layer
- Integration tests for commands
- Shell script tests for integration
- Cross-platform CI/CD

### Release Process
1. Version tagging
2. Automated builds
3. GitHub releases
4. Homebrew formula
5. AUR package

## Potential Issues & Solutions

### Issue: Directory No Longer Exists
**Solution**: Mark as invalid in list, prompt to delete

### Issue: Conflicting Aliases
**Solution**: Prompt to overwrite with confirmation

### Issue: Large Number of Bookmarks
**Solution**: Implement search/filter, pagination

### Issue: Shell Compatibility
**Solution**: Detect shell type, provide appropriate function

### Issue: Permission Errors
**Solution**: Graceful degradation, clear error messages