#!/bin/bash

# fn - Fast Navigation Tool
# One-line installer script
# Usage: curl -sSL https://raw.githubusercontent.com/user/fn/main/scripts/install-online.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/rethil/fn"
BINARY_NAME="fn"
INSTALL_DIR="/usr/local/bin"
SHELL_FUNCTION_NAME="fn"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    arm*) ARCH="arm" ;;
    *) 
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin) ;;
    *)
        echo -e "${RED}Error: Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

# Check if running as root for system-wide install
if [[ $EUID -eq 0 ]]; then
    INSTALL_DIR="/usr/local/bin"
    echo -e "${YELLOW}Installing system-wide to $INSTALL_DIR${NC}"
else
    INSTALL_DIR="$HOME/.local/bin"
    echo -e "${YELLOW}Installing to user directory $INSTALL_DIR${NC}"
    mkdir -p "$INSTALL_DIR"
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo -e "${YELLOW}Note: $INSTALL_DIR is not in your PATH${NC}"
        echo -e "${YELLOW}Consider adding this to your shell profile:${NC}"
        echo -e "${BLUE}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
    fi
fi

echo -e "${BLUE}Installing fn for $OS-$ARCH...${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# For now, since we don't have releases yet, we'll build from source
echo -e "${YELLOW}Building from source (releases coming soon)...${NC}"

# Check for required tools
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is required to build from source${NC}"
    echo -e "${YELLOW}Please install Go from https://golang.org/dl/${NC}"
    echo -e "${YELLOW}Or wait for pre-built releases to be available${NC}"
    exit 1
fi

if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git is required to clone the repository${NC}"
    exit 1
fi

# Clone repository
echo -e "${BLUE}Cloning repository...${NC}"
cd "$TMP_DIR"
git clone "$REPO_URL" fn-repo
cd fn-repo

# Build binary
echo -e "${BLUE}Building binary...${NC}"
go build -o "$BINARY_NAME" -ldflags "-s -w"

# Install binary
echo -e "${BLUE}Installing binary to $INSTALL_DIR/$BINARY_NAME...${NC}"
if [[ $EUID -eq 0 ]] || [[ -w "$INSTALL_DIR" ]]; then
    cp "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${YELLOW}Requesting sudo permissions to install to $INSTALL_DIR...${NC}"
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Detect shell
SHELL_NAME=$(basename "$SHELL")
SHELL_RC=""

case "$SHELL_NAME" in
    bash)
        if [[ -f "$HOME/.bashrc" ]]; then
            SHELL_RC="$HOME/.bashrc"
        elif [[ -f "$HOME/.bash_profile" ]]; then
            SHELL_RC="$HOME/.bash_profile"
        fi
        ;;
    zsh)
        SHELL_RC="$HOME/.zshrc"
        ;;
    fish)
        echo -e "${YELLOW}Fish shell detected. Manual setup required.${NC}"
        echo -e "${YELLOW}Add this to your Fish config:${NC}"
        echo -e "${BLUE}function fn --description 'Fast navigation with bookmarks'${NC}"
        echo -e "${BLUE}    set result (command fn navigate \$argv 2>/dev/null)${NC}"
        echo -e "${BLUE}    if test \$status -eq 0${NC}"
        echo -e "${BLUE}        cd \$result${NC}"
        echo -e "${BLUE}    else${NC}"
        echo -e "${BLUE}        command fn \$argv${NC}"
        echo -e "${BLUE}    end${NC}"
        echo -e "${BLUE}end${NC}"
        ;;
    *)
        echo -e "${YELLOW}Unknown shell: $SHELL_NAME${NC}"
        echo -e "${YELLOW}Manual setup required.${NC}"
        ;;
esac

# Install shell function for bash/zsh
if [[ -n "$SHELL_RC" ]]; then
    echo -e "${BLUE}Installing shell function to $SHELL_RC...${NC}"
    
    # Check if function already exists
    if grep -q "fn() {" "$SHELL_RC" 2>/dev/null; then
        echo -e "${YELLOW}Shell function already exists in $SHELL_RC${NC}"
    else
        cat >> "$SHELL_RC" << 'EOF'

# fn - Fast Navigation Tool
fn() {
    if [[ $# -eq 0 ]]; then
        command fn --help
        return
    fi
    
    case "$1" in
        save|list|delete|edit|cleanup|search|recent|r|path|help|completion)
            command fn "$@"
            ;;
        *)
            # Try to navigate
            local result
            result=$(command fn navigate "$1" 2>/dev/null)
            if [[ $? -eq 0 && -n "$result" ]]; then
                cd "$result"
            else
                command fn "$@"
            fi
            ;;
    esac
}
EOF
        echo -e "${GREEN}Shell function added to $SHELL_RC${NC}"
    fi
fi

# Verify installation
echo -e "${BLUE}Verifying installation...${NC}"
if command -v "$BINARY_NAME" &> /dev/null; then
    VERSION=$("$BINARY_NAME" --version 2>/dev/null || echo "latest")
    echo -e "${GREEN}✓ fn installed successfully! ${VERSION}${NC}"
    
    echo -e "${BLUE}Getting started:${NC}"
    echo -e "  ${YELLOW}fn save work${NC}      # Save current directory as 'work'"
    echo -e "  ${YELLOW}fn work${NC}           # Navigate to work directory"
    echo -e "  ${YELLOW}fn list${NC}           # List all bookmarks"
    echo -e "  ${YELLOW}fn recent${NC}         # Show recently used bookmarks"
    echo -e "  ${YELLOW}fn help${NC}           # Show all commands"
    
    if [[ -n "$SHELL_RC" ]]; then
        echo -e "${YELLOW}Restart your shell or run: source $SHELL_RC${NC}"
    fi
else
    echo -e "${RED}Error: Installation verification failed${NC}"
    exit 1
fi

echo -e "${GREEN}Installation complete! 🎉${NC}"