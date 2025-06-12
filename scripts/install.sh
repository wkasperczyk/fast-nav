#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $arch in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) print_error "Unsupported architecture: $arch"; exit 1 ;;
    esac
    
    case $os in
        linux|darwin) ;;
        *) print_error "Unsupported OS: $os"; exit 1 ;;
    esac
    
    echo "${os}_${arch}"
}

# Install binary
install_binary() {
    local platform=$1
    local install_dir="/usr/local/bin"
    
    if [[ ! -w "$install_dir" ]]; then
        print_error "Cannot write to $install_dir. Please run with sudo or add to PATH manually."
        exit 1
    fi
    
    print_status "Installing fn to $install_dir..."
    cp fn "$install_dir/fn"
    chmod +x "$install_dir/fn"
    print_status "Binary installed successfully"
}

# Add shell function
add_shell_function() {
    local shell_config=""
    local shell_name=$(basename "$SHELL")
    
    case $shell_name in
        bash)
            if [[ -f "$HOME/.bashrc" ]]; then
                shell_config="$HOME/.bashrc"
            elif [[ -f "$HOME/.bash_profile" ]]; then
                shell_config="$HOME/.bash_profile"
            fi
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            ;;
        *)
            print_warning "Unsupported shell: $shell_name. Please add the shell function manually."
            return
            ;;
    esac
    
    if [[ -z "$shell_config" ]]; then
        print_warning "Could not find shell config file. Please add the shell function manually."
        return
    fi
    
    # Check if function already exists
    if grep -q "fn()" "$shell_config"; then
        print_warning "Shell function already exists in $shell_config"
        return
    fi
    
    print_status "Adding shell function to $shell_config..."
    
    cat >> "$shell_config" << 'EOF'

# fn - Fast Navigation
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
EOF
    
    print_status "Shell function added successfully"
    print_status "Please restart your shell or run: source $shell_config"
}

# Main installation
main() {
    print_status "Installing fn - Fast Navigation Tool"
    
    # Check if binary exists
    if [[ ! -f "fn" ]]; then
        print_error "Binary 'fn' not found. Please build it first with: go build -o fn"
        exit 1
    fi
    
    # Detect platform
    local platform=$(detect_platform)
    print_status "Detected platform: $platform"
    
    # Install binary
    install_binary "$platform"
    
    # Add shell function
    add_shell_function
    
    print_status "Installation complete!"
    print_status "Usage:"
    print_status "  fn save <alias>     - Save current directory"
    print_status "  fn <alias>          - Navigate to saved directory"
    print_status "  fn list             - List all bookmarks"
    print_status "  fn delete <alias>   - Delete a bookmark"
    print_status "  fn path <alias>     - Show path without navigating"
}

main "$@"