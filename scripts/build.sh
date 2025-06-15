#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    print_status "Building for $os/$arch..."
    
    export GOOS=$os
    export GOARCH=$arch
    
    if [[ "$os" == "windows" ]]; then
        go build -o "dist/${output_name}.exe" -ldflags "-s -w"
    else
        go build -o "dist/${output_name}" -ldflags "-s -w"
    fi
    
    if [[ $? -eq 0 ]]; then
        print_status "Built: dist/${output_name}$(if [[ "$os" == "windows" ]]; then echo ".exe"; fi)"
    else
        print_error "Failed to build for $os/$arch"
        return 1
    fi
}

# Main build function
main() {
    local target_platform=""
    local build_all=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --platform)
                target_platform="$2"
                shift 2
                ;;
            --all)
                build_all=true
                shift
                ;;
            --help|-h)
                echo "Usage: $0 [--platform OS/ARCH] [--all]"
                echo ""
                echo "Options:"
                echo "  --platform OS/ARCH  Build for specific platform (e.g., linux/amd64, windows/amd64)"
                echo "  --all                Build for all supported platforms"
                echo ""
                echo "Supported platforms:"
                echo "  linux/amd64, linux/arm64, linux/386"
                echo "  darwin/amd64, darwin/arm64"
                echo "  windows/amd64, windows/386"
                exit 0
                ;;
            *)
                print_error "Unknown argument: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    # Create dist directory
    mkdir -p dist
    
    if [[ "$build_all" == "true" ]]; then
        print_status "Building for all platforms..."
        
        # Linux
        build_platform "linux" "amd64" "fn-linux-amd64"
        build_platform "linux" "arm64" "fn-linux-arm64"
        build_platform "linux" "386" "fn-linux-386"
        
        # macOS
        build_platform "darwin" "amd64" "fn-darwin-amd64"
        build_platform "darwin" "arm64" "fn-darwin-arm64"
        
        # Windows
        build_platform "windows" "amd64" "fn-windows-amd64"
        build_platform "windows" "386" "fn-windows-386"
        
        print_status "All builds completed successfully!"
        print_status "Binaries are in the dist/ directory"
        
    elif [[ -n "$target_platform" ]]; then
        IFS='/' read -r os arch <<< "$target_platform"
        
        if [[ -z "$os" ]] || [[ -z "$arch" ]]; then
            print_error "Invalid platform format. Use OS/ARCH (e.g., linux/amd64)"
            exit 1
        fi
        
        build_platform "$os" "$arch" "fn-${os}-${arch}"
        
    else
        # Build for current platform
        local current_os=$(go env GOOS)
        local current_arch=$(go env GOARCH)
        
        print_status "Building for current platform ($current_os/$current_arch)..."
        
        if [[ "$current_os" == "windows" ]]; then
            go build -o "fast-nav.exe" -ldflags "-s -w"
            print_status "Built: fast-nav.exe"
        else
            go build -o "fast-nav" -ldflags "-s -w"
            print_status "Built: fast-nav"
        fi
    fi
}

main "$@"