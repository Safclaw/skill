# Skill Installer - One-line installation for Windows
# Usage: Invoke-WebRequest -UseBasicParsing https://raw.githubusercontent.com/Safclaw/skill/main/scripts/install.ps1 | Invoke-Expression

param(
    [string]$Version = "latest",
    [switch]$Help
)

$ErrorActionPreference = "Stop"

# Helper functions
function Write-Info {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

# Show help
if ($Help) {
    Write-Host "Usage: Invoke-WebRequest -UseBasicParsing https://raw.githubusercontent.com/Safclaw/skill/main/scripts/install.ps1 | Invoke-Expression"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version VERSION  Install specific version"
    Write-Host "  -Help             Show this help message"
    exit 0
}

Write-Host "🚀 Skill Installer"
Write-Host "=================="
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
Write-Info "Detected: windows/$Arch"

# Determine install location
$InstallDir = "$env:ProgramFiles\skill"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
    Write-Info "Created install directory: $InstallDir"
}

# Build download URL
$Repo = "Safclaw/skill"
$BinaryName = "skill_windows_${Arch}.exe"
$DownloadUrl = "https://github.com/${Repo}/releases/latest/download/${BinaryName}"

if ($Version -ne "latest") {
    $DownloadUrl = "https://github.com/${Repo}/releases/download/${Version}/${BinaryName}"
}

# Download and install
Write-Info "Downloading skill ${Version}..."

try {
    $TempFile = [System.IO.Path]::GetTempFileName()
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
    
    # Move to install directory
    $TargetPath = Join-Path $InstallDir "skill.exe"
    Move-Item -Path $TempFile -Destination $TargetPath -Force
    
    Write-Info "Successfully installed skill to $TargetPath"
} catch {
    Write-Error-Custom "Failed to download: $_"
    exit 1
}

# Add to PATH if not already present
$CurrentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentUserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$CurrentUserPath;$InstallDir", "User")
    Write-Info "Added $InstallDir to PATH"
    Write-Warn "Please restart your terminal or run: `$env:Path = [Environment]::GetEnvironmentVariable(`"Path`",`"Machine`") + `";`" + [Environment]::GetEnvironmentVariable(`"Path`",`"User`")"
}

# Verify installation
try {
    & "$TargetPath" --version
    Write-Info "Installation verified!"
    Write-Host ""
    Write-Host "You can now use:"
    Write-Host "  skill add github.com/Safclaw/skills/read-json"
    Write-Host "  skill --help"
} catch {
    Write-Error-Custom "Installation verification failed"
    exit 1
}

Write-Info "Installation complete!"
