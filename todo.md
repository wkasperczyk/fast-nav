# TODO - Fast Navigation Tool

This file tracks all planned improvements and tasks for the fn application. Check this file regularly to see what needs to be done next.

## Phase 1: Core Completion (HIGH PRIORITY)

### Missing Core Features
- [x] **Implement edit command** - Update existing alias to current directory (mentioned in root.go but not implemented)
- [x] **Add cleanup command** - Remove dead/missing bookmarks automatically
- [x] **Add search/filter functionality** - Find bookmarks by path or alias patterns
- [x] **Add tab completion** - Shell completion for aliases

### Testing Infrastructure (CRITICAL)
- [x] **Create storage_test.go** - Unit tests for storage layer
- [ ] **Add CLI command tests** - Test all cobra commands
- [ ] **Add integration tests** - End-to-end testing
- [ ] **Set up test coverage reporting** - Track test coverage
- [ ] **Add benchmark tests** - Performance testing for storage operations

### Error Handling & Reliability
- [ ] **Improve error messages** - More user-friendly error output
- [ ] **Add logging system** - Debug logging for troubleshooting
- [ ] **Handle edge cases** - Empty aliases, special characters, etc.

## Phase 2: User Experience

### Enhanced Navigation
- [ ] **Fuzzy matching for navigation** - Partial alias matching
- [ ] **Recently used shortcuts** - Quick access to frequent bookmarks
- [ ] **Smart suggestions** - Suggest similar aliases on typos

### Installation & Setup
- [ ] **One-line install script** - curl | bash style installer
- [ ] **Uninstall functionality** - Clean removal of binary and shell function
- [ ] **Auto-detection of shell type** - Better shell configuration detection
- [ ] **Windows support** - PowerShell integration

### Shell Integration
- [ ] **Bash completion script** - Tab completion for bash
- [ ] **Zsh completion script** - Tab completion for zsh
- [ ] **Fish shell support** - Support for fish shell users

## Phase 3: Advanced Features

### Data Management
- [ ] **Import/Export functionality** - Backup bookmarks to file
- [ ] **Sync between machines** - Share bookmarks across systems
- [ ] **Export to different formats** - JSON, CSV, plain text export

### Enhanced Listing & Search
- [ ] **Sort options** - Sort by usage, date, name, path
- [ ] **Tree view** - Hierarchical view of paths
- [ ] **Search by path patterns** - Find bookmarks containing path segments
- [ ] **Usage statistics** - Detailed usage reports

### Configuration System
- [ ] **Custom storage location** - Allow user to specify bookmark file location
- [ ] **Alias validation rules** - Configurable alias format rules
- [ ] **Output formatting options** - Customizable list output
- [ ] **Color scheme configuration** - User-defined color preferences

## Phase 4: Production Polish

### Build & Release
- [ ] **Create Makefile** - Standardized build process
- [ ] **Cross-platform builds** - Linux, macOS, Windows binaries
- [ ] **GitHub Actions CI/CD** - Automated testing and releases
- [ ] **Automated releases** - Tag-based release automation
- [ ] **Release binaries** - Pre-built binaries for each platform

### Distribution
- [ ] **Homebrew formula** - macOS package manager
- [ ] **APT/DEB packages** - Debian/Ubuntu packages
- [ ] **RPM packages** - RedHat/Fedora packages
- [ ] **AUR package** - Arch Linux package
- [ ] **Docker image** - Containerized version
- [ ] **Snap package** - Universal Linux package

### Documentation
- [ ] **Man pages** - Standard Unix documentation
- [ ] **Comprehensive examples** - Usage examples and tutorials
- [ ] **Troubleshooting guide** - Common issues and solutions
- [ ] **Migration guide** - From other bookmark tools
- [ ] **API documentation** - For developers extending the tool

## Current Priority (Immediate Next Steps)

1. **Implement edit command** (quick win - already mentioned in help)
2. **Set up basic testing infrastructure** (critical for reliability)
3. **Create Makefile** (consistent builds)
4. **Implement cleanup command** (addresses user pain point with missing directories)

## Adding New Tasks

When new tasks or improvements are identified:
1. Add them to the appropriate phase in this file
2. Include priority level and brief description
3. Link related tasks together
4. Update completion status as work progresses

## Notes

- Mark completed tasks with [x] 
- Add new phases as needed
- Keep this file updated as the primary source of truth for planned work
- Reference this file in commit messages when completing tasks