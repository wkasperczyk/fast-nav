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
        i686|i386) arch="386" ;;
        *) print_error "Unsupported architecture: $arch"; exit 1 ;;
    esac
    
    case $os in
        linux|darwin) ;;
        *) print_error "Unsupported OS: $os. For Windows, use install.ps1"; exit 1 ;;
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

# Detect shell type and configuration file
detect_shell() {
    local shell_name=$(basename "$SHELL" 2>/dev/null || echo "")
    local shell_config=""
    local shell_function_type="bash"
    
    # If SHELL is not set, try to detect from process tree
    if [[ -z "$shell_name" ]]; then
        shell_name=$(ps -p $$ -o comm= 2>/dev/null | sed 's/^-//' || echo "sh")
    fi
    
    print_status "Detected shell: $shell_name"
    
    case $shell_name in
        bash)
            # macOS prefers .bash_profile, Linux prefers .bashrc
            if [[ "$(uname -s)" == "Darwin" ]]; then
                if [[ -f "$HOME/.bash_profile" ]]; then
                    shell_config="$HOME/.bash_profile"
                else
                    shell_config="$HOME/.bashrc"
                fi
            else
                if [[ -f "$HOME/.bashrc" ]]; then
                    shell_config="$HOME/.bashrc"
                else
                    shell_config="$HOME/.bash_profile"
                fi
            fi
            shell_function_type="bash"
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            shell_function_type="zsh"
            ;;
        fish)
            # Create fish config directory if it doesn't exist
            mkdir -p "$HOME/.config/fish"
            shell_config="$HOME/.config/fish/config.fish"
            shell_function_type="fish"
            ;;
        sh)
            # POSIX shell - try common config files
            if [[ -f "$HOME/.profile" ]]; then
                shell_config="$HOME/.profile"
            else
                shell_config="$HOME/.shellrc"
            fi
            shell_function_type="posix"
            ;;
        ksh|mksh)
            if [[ -f "$HOME/.kshrc" ]]; then
                shell_config="$HOME/.kshrc"
            else
                shell_config="$HOME/.profile"
            fi
            shell_function_type="bash"  # ksh is mostly bash-compatible
            ;;
        tcsh|csh)
            shell_config="$HOME/.cshrc"
            shell_function_type="csh"
            ;;
        *)
            print_warning "Unsupported shell: $shell_name"
            print_status "Supported shells: bash, zsh, fish, sh, ksh, mksh, tcsh, csh"
            print_status "Please add the shell function manually or switch to a supported shell."
            return 1
            ;;
    esac
    
    echo "$shell_config|$shell_function_type"
}

# Generate shell function based on shell type
generate_shell_function() {
    local function_type="$1"
    
    case $function_type in
        bash|zsh|posix)
            cat << 'EOF'

# fn - Fast Navigation
fn() {
    if [[ "$1" == "save" ]] || [[ "$1" == "list" ]] || [[ "$1" == "delete" ]] || [[ "$1" == "path" ]] || [[ "$1" == "edit" ]] || [[ "$1" == "cleanup" ]]; then
        command fn "$@"
    else
        local dir=$(command fn navigate "$@")
        if [[ -n "$dir" ]]; then
            cd "$dir"
        fi
    fi
}
EOF
            ;;
        fish)
            cat << 'EOF'

# fn - Fast Navigation
function fn
    if contains $argv[1] save list delete path edit cleanup
        command fn $argv
    else
        set dir (command fn navigate $argv)
        if test -n "$dir"
            cd "$dir"
        end
    end
end
EOF
            ;;
        csh)
            cat << 'EOF'

# fn - Fast Navigation
alias fn 'set args = (\!*); if ("$args[1]" == "save" || "$args[1]" == "list" || "$args[1]" == "delete" || "$args[1]" == "path" || "$args[1]" == "edit" || "$args[1]" == "cleanup") then; command fn $args; else; set dir = `command fn navigate $args`; if ("$dir" != "") cd "$dir"; endif'
EOF
            ;;
    esac
}

# Add shell function
add_shell_function() {
    local shell_info=$(detect_shell)
    if [[ $? -ne 0 ]]; then
        return 1
    fi
    
    local shell_config=$(echo "$shell_info" | cut -d'|' -f1)
    local function_type=$(echo "$shell_info" | cut -d'|' -f2)
    
    # Create config file directory if it doesn't exist
    mkdir -p "$(dirname "$shell_config")"
    
    # Check if function already exists (multiple patterns for different shells)
    local patterns=("fn()" "function fn" "alias fn")
    local has_function=false
    
    if [[ -f "$shell_config" ]]; then
        for pattern in "${patterns[@]}"; do
            if grep -q "$pattern" "$shell_config"; then
                has_function=true
                break
            fi
        done
    fi
    
    if [[ "$has_function" == "true" ]]; then
        print_warning "Shell function already exists in $shell_config"
        return 0
    fi
    
    print_status "Adding shell function to $shell_config..."
    
    # Generate and append the appropriate shell function
    generate_shell_function "$function_type" >> "$shell_config"
    
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