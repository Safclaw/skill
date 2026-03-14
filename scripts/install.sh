#!/bin/bash
# Skill Installer - One-line installation for Mac and Linux
# Usage: curl -fsSL https://raw.githubusercontent.com/Safclaw/skill/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="Safclaw/skill"
INSTALL_DIR="/usr/local/bin"
SKILL_BIN="skill"

# Helper functions
log_info() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

# Detect OS and architecture
detect_os() {
    case "$(uname -s)" in
        Darwin)
            echo "darwin"
            ;;
        Linux)
            echo "linux"
            ;;
        *)
            log_error "Unsupported OS: $(uname -s)"
            exit 1
            ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Check prerequisites
check_prerequisites() {
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed."
        exit 1
    fi
    
    if [ "$(id -u)" -ne 0 ]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Download and install
install_skill() {
    local os="$1"
    local arch="$2"
    local version="$3"
    
    local binary_name="skill_${os}_${arch}"
    local download_url="https://github.com/${REPO}/releases/latest/download/${binary_name}"
    
    if [ "$os" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi
    
    log_info "Downloading ${SKILL_BIN} ${version:-latest} for ${os}/${arch}..."
    
    # Create temporary directory
    local tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT
    
    # Download binary
    if ! curl -sL -o "${tmp_dir}/${SKILL_BIN}" "${download_url}"; then
        log_error "Failed to download from ${download_url}"
        exit 1
    fi
    
    # Make executable
    chmod +x "${tmp_dir}/${SKILL_BIN}"
    
    # Install to target directory
    if ! mv "${tmp_dir}/${SKILL_BIN}" "${INSTALL_DIR}/${SKILL_BIN}"; then
        log_error "Failed to install to ${INSTALL_DIR}"
        exit 1
    fi
    
    log_info "Successfully installed ${SKILL_BIN} to ${INSTALL_DIR}"
}

# Verify installation
verify_installation() {
    if command -v ${SKILL_BIN} &> /dev/null; then
        local version=$(${SKILL_BIN} --version 2>&1 || echo "unknown")
        log_info "Installation verified! ${SKILL_BIN} ${version}"
        echo ""
        echo "You can now use:"
        echo "  ${SKILL_BIN} add github.com/Safclaw/skills/read-json"
        echo "  ${SKILL_BIN} --help"
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

# Main execution
main() {
    echo "🚀 Skill Installer"
    echo "=================="
    echo ""
    
    check_prerequisites
    
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    log_info "Detected: ${os}/${arch}"
    
    install_skill "${os}" "${arch}"
    verify_installation
    
    log_info "Installation complete!"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: curl -fsSL https://raw.githubusercontent.com/Safclaw/skill/main/scripts/install.sh | bash"
            echo ""
            echo "Options:"
            echo "  -v, --version VERSION  Install specific version"
            echo "  -h, --help             Show this help message"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

main
