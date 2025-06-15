#!/usr/bin/env bash

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
        *) print_error "Unsupported architecture: $arch"; return 1 ;;
    esac
    
    case $os in
        linux|darwin) ;;
        *) print_error "Unsupported OS: $os. For Windows, use install.ps1"; return 1 ;;
    esac
    
    echo "${os}_${arch}"
}

# Find suitable installation directory
find_install_dir() {
    local dirs=("/usr/local/bin" "$HOME/.local/bin" "$HOME/bin")
    
    for dir in "${dirs[@]}"; do
        if [[ -d "$dir" ]] && [[ -w "$dir" ]]; then
            echo "$dir"
            return 0
        fi
    done
    
    # Create ~/.local/bin if it doesn't exist
    if mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        echo "$HOME/.local/bin"
        return 0
    fi
    
    return 1
}

# Install binary
install_binary() {
    local platform=$1
    local install_dir=$(find_install_dir)
    
    if [[ -z "$install_dir" ]]; then
        print_error "Cannot find a writable installation directory"
        print_status "Try running with sudo or create $HOME/.local/bin"
        return 1
    fi
    
    print_status "Installing fast-nav to $install_dir..."
    
    if ! cp fast-nav "$install_dir/fast-nav" 2>/dev/null; then
        print_error "Failed to copy binary. Check permissions."
        return 1
    fi
    
    if ! chmod +x "$install_dir/fast-nav" 2>/dev/null; then
        print_error "Failed to make binary executable"
        return 1
    fi
    
    # Check if directory is in PATH
    if ! echo "$PATH" | grep -q "$install_dir"; then
        print_warning "$install_dir is not in your PATH"
        print_status "Add this line to your shell config: export PATH=\"$install_dir:\$PATH\""
    fi
    
    print_status "Binary installed successfully"
    return 0
}

# Detect shell type and configuration file
detect_shell() {
    local shell_name=""
    local shell_config=""
    local shell_function_type="bash"

    # Try multiple methods to detect shell
    if [[ -n "$ZSH_VERSION" ]]; then
        shell_name="zsh"
    elif [[ -n "$BASH_VERSION" ]]; then
        shell_name="bash"
    elif [[ -n "$FISH_VERSION" ]]; then
        shell_name="fish"
    elif [[ -n "$KSH_VERSION" ]]; then
        shell_name="ksh"
    else
        # Try to get shell from parent process or $SHELL
        if command -v ps >/dev/null 2>&1; then
            shell_name=$(ps -p $$ -o comm= 2>/dev/null | sed 's/^-//' | sed 's/.*\///')
            if [[ -z "$shell_name" || "$shell_name" == "sh" ]]; then
                local ppid=$(ps -p $$ -o ppid= 2>/dev/null | tr -d ' ')
                if [[ -n "$ppid" ]]; then
                    shell_name=$(ps -p $ppid -o comm= 2>/dev/null | sed 's/^-//' | sed 's/.*\///')
                fi
            fi
        fi
        
        if [[ -z "$shell_name" || "$shell_name" == "sh" ]]; then
            shell_name=$(basename "${SHELL:-sh}")
        fi
    fi

    print_status "Detected shell: $shell_name"

    case $shell_name in
        bash)
            # Check for existing config files in order of preference
            for config in "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile"; do
                if [[ -f "$config" ]]; then
                    shell_config="$config"
                    break
                fi
            done
            # Create .bashrc if no config exists
            [[ -z "$shell_config" ]] && shell_config="$HOME/.bashrc"
            shell_function_type="bash"
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            shell_function_type="bash"  # zsh is bash-compatible for our purposes
            ;;
        fish)
            mkdir -p "$HOME/.config/fish/functions" 2>/dev/null
            shell_config="$HOME/.config/fish/functions/fn.fish"
            shell_function_type="fish"
            ;;
        sh|dash|ash)
            shell_config="$HOME/.profile"
            shell_function_type="posix"
            ;;
        ksh|mksh|pdksh)
            for config in "$HOME/.kshrc" "$HOME/.profile"; do
                if [[ -f "$config" ]]; then
                    shell_config="$config"
                    break
                fi
            done
            [[ -z "$shell_config" ]] && shell_config="$HOME/.kshrc"
            shell_function_type="bash"
            ;;
        tcsh|csh)
            shell_config="$HOME/.cshrc"
            shell_function_type="csh"
            ;;
        *)
            print_warning "Unknown shell: $shell_name"
            print_status "You'll need to add the function manually"
            return 1
            ;;
    esac

    echo "$shell_config|$shell_function_type"
}

# Generate shell function based on shell type
generate_shell_function() {
    local function_type="$1"
    
    case $function_type in
        bash|posix)
            cat << 'EOF'

# fn - Fast Navigation
fn() {
    case "$1" in
        save|list|delete|path|edit|cleanup)
            command fast-nav "$@"
            ;;
        *)
            local dir
            dir=$(command fast-nav navigate "$@" 2>/dev/null)
            if [ -n "$dir" ] && [ -d "$dir" ]; then
                cd "$dir" || return 1
            fi
            ;;
    esac
}
EOF
            ;;
        fish)
            # For fish, we create a complete function file
            cat << 'EOF'
# fn - Fast Navigation
function fn
    switch $argv[1]
        case save list delete path edit cleanup
            command fast-nav $argv
        case '*'
            set -l dir (command fast-nav navigate $argv 2>/dev/null)
            if test -n "$dir" -a -d "$dir"
                cd "$dir"
            end
    end
end
EOF
            ;;
        csh)
            cat << 'EOF'

# fn - Fast Navigation
alias fn 'if ("\!:1" == "save" || "\!:1" == "list" || "\!:1" == "delete" || "\!:1" == "path" || "\!:1" == "edit" || "\!:1" == "cleanup") then \\
    command fast-nav \!* \\
else \\
    set fn_dir = `command fast-nav navigate \!* |& grep -v "^$"` \\
    if ( "$fn_dir" != "" && -d "$fn_dir" ) cd "$fn_dir" \\
endif'
EOF
            ;;
    esac
}

# Check if function already exists
function_exists() {
    local shell_config="$1"
    local function_type="$2"
    
    if [[ ! -f "$shell_config" ]]; then
        return 1
    fi
    
    case $function_type in
        fish)
            # Fish function files are self-contained
            return 0
            ;;
        csh)
            grep -q "alias fn" "$shell_config" 2>/dev/null
            ;;
        *)
            grep -E "^[[:space:]]*(function[[:space:]]+)?fn[[:space:]]*\(\)" "$shell_config" 2>/dev/null
            ;;
    esac
}

# Add shell function
add_shell_function() {
    local shell_info
    shell_info=$(detect_shell)
    if [[ $? -ne 0 ]]; then
        return 1
    fi
    
    local shell_config=$(echo "$shell_info" | cut -d'|' -f1)
    local function_type=$(echo "$shell_info" | cut -d'|' -f2)
    
    # Create config file directory if needed
    local config_dir=$(dirname "$shell_config")
    if [[ ! -d "$config_dir" ]]; then
        mkdir -p "$config_dir" 2>/dev/null || {
            print_error "Cannot create directory: $config_dir"
            return 1
        }
    fi
    
    # Check if function already exists
    if function_exists "$shell_config" "$function_type"; then
        print_warning "Shell function already exists in $shell_config"
        return 0
    fi
    
    print_status "Adding shell function to $shell_config..."
    
    # For fish, we write directly to the function file
    if [[ "$function_type" == "fish" ]]; then
        generate_shell_function "$function_type" > "$shell_config"
    else
        # For other shells, append to config
        {
            echo ""  # Add newline before our content
            generate_shell_function "$function_type"
        } >> "$shell_config"
    fi
    
    if [[ $? -eq 0 ]]; then
        print_status "Shell function added successfully"
        
        # Provide appropriate reload command
        case $function_type in
            fish)
                print_status "Please restart your shell or run: source $shell_config"
                ;;
            csh)
                print_status "Please restart your shell or run: source $shell_config"
                ;;
            *)
                print_status "Please restart your shell or run: source $shell_config"
                ;;
        esac
    else
        print_error "Failed to add shell function"
        return 1
    fi
    
    return 0
}

# Main installation
main() {
    print_status "Installing fn - Fast Navigation Tool"
    
    # Check if binary exists
    if [[ ! -f "fast-nav" ]]; then
        print_error "Binary 'fast-nav' not found in current directory"
        print_status "Please build it first with: go build -o fast-nav"
        return 1
    fi
    
    # Make sure binary is executable
    if [[ ! -x "fast-nav" ]]; then
        chmod +x fast-nav 2>/dev/null || {
            print_error "Cannot make binary executable"
            return 1
        }
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    if [[ $? -ne 0 ]]; then
        return 1
    fi
    print_status "Detected platform: $platform"
    
    # Install binary
    if ! install_binary "$platform"; then
        return 1
    fi
    
    # Add shell function
    add_shell_function
    
    echo ""
    print_status "Installation complete!"
    print_status ""
    print_status "Usage:"
    print_status "  fn save <alias>     - Save current directory"
    print_status "  fn <alias>          - Navigate to saved directory"
    print_status "  fn list             - List all bookmarks"
    print_status "  fn delete <alias>   - Delete a bookmark"
    print_status "  fn path <alias>     - Show path without navigating"
    print_status "  fn edit             - Edit bookmarks file"
    print_status "  fn cleanup          - Remove invalid bookmarks"
    
    return 0
}

# Run main function
main "$@"
exit $?