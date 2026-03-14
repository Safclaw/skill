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
REPO="safeclaw/skill"
INSTALL_DIR="/usr/local/bin"
SKILL_BIN="skill"
USER_INSTALL_DIR="$HOME/.local/bin"

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
}

# Detect if running as root
is_root() {
    [ "$(id -u)" -eq 0 ]
}

# Choose install directory based on permissions
choose_install_dir() {
    if is_root || [ -w "$INSTALL_DIR" ]; then
        echo "$INSTALL_DIR"
    else
        # Check if user has write permission to INSTALL_DIR
        if [ ! -d "$INSTALL_DIR" ]; then
            # Try to create with sudo
            log_warn "Cannot write to $INSTALL_DIR, will use $USER_INSTALL_DIR"
            echo "$USER_INSTALL_DIR"
        elif [ -w "$INSTALL_DIR" ]; then
            echo "$INSTALL_DIR"
        else
            log_warn "Cannot write to $INSTALL_DIR, will use $USER_INSTALL_DIR"
            echo "$USER_INSTALL_DIR"
        fi
    fi
}

# Download and install
install_skill() {
    local os="$1"
    local arch="$2"
    local version="$3"
    local target_dir="$4"
    
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
    
    # Check if downloaded file is valid (not an HTML error page)
    if [ -f "${tmp_dir}/${SKILL_BIN}" ]; then
        # Check first few bytes for ELF or Mach-O magic numbers
        local first_bytes=$(head -c 4 "${tmp_dir}/${SKILL_BIN}" 2>/dev/null || echo "")
        if [[ "$first_bytes" != $'\x7fELF' ]] && [[ "$first_bytes" != $'\xcf\xfa\xed\xfe' ]] && [[ "$first_bytes" != $'\xfe\xed\xfa\xcf' ]]; then
            # Check if it's a text file (likely an error message)
            if file "${tmp_dir}/${SKILL_BIN}" | grep -q "text\|HTML"; then
                log_error "Downloaded file is not a valid binary. The release may not exist."
                log_error "Please check: ${download_url}"
                exit 1
            fi
        fi
    fi
    
    # Make executable and verify it's a valid binary
    chmod +x "${tmp_dir}/${SKILL_BIN}"
    
    # Verify the binary is executable
    if ! "${tmp_dir}/${SKILL_BIN}" --version &>/dev/null; then
        log_warn "Binary verification failed, but continuing with installation..."
    fi
    
    # Create target directory if needed
    if [ ! -d "$target_dir" ]; then
        mkdir -p "$target_dir"
    fi
    
    # Install to target directory
    if ! mv "${tmp_dir}/${SKILL_BIN}" "${target_dir}/${SKILL_BIN}"; then
        # If failed and trying system dir, try with sudo
        if [ "$target_dir" = "$INSTALL_DIR" ]; then
            log_info "Attempting installation with sudo..."
            if sudo mv "${tmp_dir}/${SKILL_BIN}" "${INSTALL_DIR}/${SKILL_BIN}"; then
                log_info "Successfully installed ${SKILL_BIN} to ${INSTALL_DIR}"
                return 0
            fi
        fi
        log_error "Failed to install to ${target_dir}"
        exit 1
    fi
    
    log_info "Successfully installed ${SKILL_BIN} to ${target_dir}"
}

# Verify installation
verify_installation() {
    # Refresh PATH if using user directory
    if [[ "$PATH" != *"$USER_INSTALL_DIR"* ]]; then
        export PATH="$USER_INSTALL_DIR:$PATH"
    fi
    
    if command -v ${SKILL_BIN} &> /dev/null; then
        local version=$(${SKILL_BIN} --version 2>&1 || echo "unknown")
        log_info "Installation verified! ${SKILL_BIN} ${version}"
        echo ""
        echo "You can now use:"
        echo "  ${SKILL_BIN} add github.com/safeclaw/skills/read-json"
        echo "  ${SKILL_BIN} --help"
    else
        log_error "Installation verification failed"
        if [ -f "$USER_INSTALL_DIR/${SKILL_BIN}" ]; then
            echo ""
            log_warn "${SKILL_BIN} is installed but not in PATH"
            echo "Add this to your ~/.bashrc or ~/.zshrc:"
            echo "  export PATH=\$HOME/.local/bin:\$PATH"
            echo "Or run: ${USER_INSTALL_DIR}/${SKILL_BIN} --help"
        fi
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
    
    # Determine install directory
    local install_dir=$(choose_install_dir)
    if [ "$install_dir" = "$USER_INSTALL_DIR" ]; then
        log_info "Installing to user directory: ${install_dir}"
    else
        log_info "Installing to system directory: ${install_dir}"
    fi
    
    install_skill "${os}" "${arch}" "${VERSION:-}" "${install_dir}"
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
