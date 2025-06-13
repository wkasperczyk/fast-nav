# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `fn` - a Go-based command-line tool for fast directory navigation using bookmarks. Users save directories with short aliases and navigate to them quickly. The tool works by outputting directory paths that are consumed by a shell function wrapper to change the working directory.

## Commands

```bash
# Build the project
go build -o fn

# Install (requires pre-built binary)
sudo ./scripts/install.sh

# Manual testing during development
./fn save test
./fn list
./fn delete test
./fn path test
```

## Architecture

### Core Components

- **cmd/**: Cobra CLI commands structure
  - `root.go`: Main command setup and subcommand registration
  - `save.go`, `list.go`, `navigate.go`, `delete.go`, `path.go`: Individual command implementations
- **internal/storage/**: Bookmark persistence layer
  - `storage.go`: JSON-based storage with `~/.fn/bookmarks.json`
  - Handles bookmark CRUD operations and usage tracking
- **main.go**: Entry point that delegates to cmd.Execute()

### Navigation Architecture

The challenge: A binary cannot change the parent shell's working directory directly.

**Solution**: The Go binary outputs the target directory path, and a shell function wrapper evaluates this output to perform the actual `cd` command.

- Direct commands (`save`, `list`, `delete`, `path`) are passed through to the binary
- Navigation commands output the target path, consumed by the shell function
- Shell function is auto-installed in `~/.bashrc`/`~/.zshrc` by the install script

### Data Structure

Bookmarks stored in `~/.fn/bookmarks.json`:
```json
{
  "version": "1.0", 
  "bookmarks": {
    "alias": {
      "path": "/full/path",
      "created": "timestamp",
      "used_count": 42,
      "last_used": "timestamp"
    }
  }
}
```

## Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/fatih/color`: Terminal colors
- `github.com/mitchellh/go-homedir`: Home directory detection
- `github.com/AlecAivazis/survey/v2`: Interactive prompts (used in some commands)

Requires Go 1.21+ and Unix-like systems (Linux/macOS).

## Task Management

**IMPORTANT**: Always check `todo.md` before starting any work. This file contains the comprehensive roadmap and prioritized task list for the project.

- `todo.md` is the authoritative source for all planned improvements and next steps
- Add any newly identified tasks or improvements to `todo.md` 
- Mark completed tasks as done in `todo.md`
- Reference `todo.md` task items in commit messages when completing work

When planning new features or fixes, update `todo.md` first to maintain a clear development roadmap.