# fn - Fast Navigation

A command-line tool that allows you to bookmark directories and quickly navigate to them using short aliases.

## Quick Start

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
```

## Installation

### Build from source

```bash
# Clone the repository
git clone https://github.com/rethil/fn.git
cd fn

# Build the binary
go build -o fn

# Install (requires sudo for /usr/local/bin)
sudo ./scripts/install.sh
```

### Manual Installation

```bash
# Build
go build -o fn

# Copy binary to PATH
sudo cp fn /usr/local/bin/

# Add shell function to your shell config (~/.bashrc, ~/.zshrc, etc.)
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

# Reload your shell
source ~/.bashrc  # or ~/.zshrc
```

## Commands

- **`fn save <alias>`** - Save current directory with an alias
- **`fn <alias>`** - Navigate to saved directory
- **`fn list`** - List all saved aliases
- **`fn delete <alias>`** - Remove a saved alias
- **`fn path <alias>`** - Print path without navigating

## How it works

The challenge with navigation tools is that a binary cannot directly change the parent shell's working directory. This tool solves it by:

1. The GO binary outputs the target directory path
2. A shell function wrapper evaluates the output and changes directory
3. Direct commands (save, list, delete, path) are passed through to the binary

## Configuration

Bookmarks are stored in `~/.fn/bookmarks.json` with the following structure:

```json
{
  "version": "1.0",
  "bookmarks": {
    "myapp": {
      "path": "/home/user/projects/myapp",
      "created": "2024-01-15T10:30:00Z",
      "used_count": 42,
      "last_used": "2024-01-20T15:45:00Z"
    }
  }
}
```

## Requirements

- Go 1.21 or later
- Unix-like system (Linux, macOS)
- Bash or Zsh shell

## License

MIT License