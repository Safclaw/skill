# skill

A Go Module-inspired decentralized skill management tool for SafeClaw.
## Install

### Quick Install (Recommended)

**Mac & Linux:**
```bash
# Non-root installation (installs to ~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/safeclaw/skill/main/scripts/install.sh | bash

# System-wide installation (requires sudo)
curl -fsSL https://raw.githubusercontent.com/safeclaw/skill/main/scripts/install.sh | sudo bash
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -UseBasicParsing https://raw.githubusercontent.com/safeclaw/skill/main/scripts/install.ps1 | Invoke-Expression
```

<details>
<summary><strong>Manual Installation</strong></summary>

#### Mac
```bash
# Option 1: Using Homebrew (if available)
brew install skill

# Option 2: Manual download (non-root)
curl -LO https://github.com/safeclaw/skill/releases/latest/download/skill_darwin_amd64
chmod +x skill_darwin_amd64
mkdir -p $HOME/.local/bin
mv skill_darwin_amd64 $HOME/.local/bin/skill
export PATH="$HOME/.local/bin:$PATH"
skill --version

# Option 3: Manual download (system-wide, requires sudo)
curl -LO https://github.com/safeclaw/skill/releases/latest/download/skill_darwin_amd64
chmod +x skill_darwin_amd64
sudo mv skill_darwin_amd64 /usr/local/bin/skill
skill --version
```

#### Linux
```bash
# Download binary (non-root)
curl -LO https://github.com/safeclaw/skill/releases/latest/download/skill_linux_amd64
chmod +x skill_linux_amd64
mkdir -p $HOME/.local/bin
mv skill_linux_amd64 $HOME/.local/bin/skill
export PATH="$HOME/.local/bin:$PATH"
skill --version

# Or system-wide (requires sudo)
curl -LO https://github.com/safeclaw/skill/releases/latest/download/skill_linux_amd64
chmod +x skill_linux_amd64
sudo mv skill_linux_amd64 /usr/local/bin/skill
skill --version
```

#### Windows
```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/safeclaw/skill/releases/latest/download/skill_windows_amd64.exe" -OutFile "$env:USERPROFILE\skill.exe"
# Add to PATH manually or move to existing PATH directory
```

</details>

## Features

- **Decentralized Management**: Download skills from any Git repository (GitHub, GitLab, Gitee, etc.)
- **Go Module Style**: Uses `{host}/{namespace}/{name}@{version}` format
- **Three Installation Modes**: Global, workspace, or custom directory
- **Hook System**: Secure lifecycle hooks with checksum verification
- **Skill Manifest**: `.skills.yaml` tracks installed skills with signatures

## Commands

### 1. Install skill
```bash
# Install to workspace
skill add github.com/safeclaw/skills/read-json -w

# Install globally
skill add -g github.com/safeclaw/skills/read-json

# Install to specified directory
skill add github.com/safeclaw/skills/read-json --workspace-path ~/.opencalw/workspace
```

### 2. Uninstall skill
```bash
# Uninstall from workspace
skill remove github.com/safeclaw/skills/read-json -w

# Uninstall global installation
skill remove -g github.com/safeclaw/skills/read-json

# Uninstall from specified directory
skill remove github.com/safeclaw/skills/read-json --workspace-path ~/.opencalw/workspace
```

### 3. Initialize a directory as a skill
```bash
# skill init {skillName} [--template github.com/safeclaw/skill/empty]
skill init github.com/xxx/xxSkill
```

After initialization, the directory structure is as follows:
```bash
├── scripts
│   ├── setup.sh # Script executed after installation
│   ├── setup.ps1 # Script executed after installation (Windows)
│   ├── unsetup.sh # Script executed before uninstallation
│   └── unsetup.ps1 # Script executed before uninstallation (Windows)
├── skill.md # Skill entry file
└── skill.yaml # Skill metadata
```

### 4. List installed skills
```bash
# List global skills
skill list -g

# List workspace skills
skill list -w

# List skills in specified directory
skill list --workspace-path /path/to/dir
```

### 5. Show skill information
```bash
skill info github.com/safeclaw/skills/read-json
```

### 6. Manage cache
```bash
# Clean cache
skill cache clean

# Clean all cache
skill cache clean --all

# Verify cache integrity
skill cache verify
```

## Installation Directory Structure

```
{install_dir}/
├── .skills.yaml          # Skill manifest
└── reps/                 # Skill storage
    └── github.com/       # Organized by host
        └── Safclaw/      # Organized by namespace
            └── skills/   # Organized by type
                ├── json/
                │   ├── skill.yaml
                │   ├── skill.md
                │   └── scripts/
                └── xlsx/
                    └── ...
```

## .skills.yaml Format

```yaml
# skill reps
---
skills:
  - name: "Commit Helper"
    dir: "github.com/commit-helper/commit-helper"
    version: v1.0.1
    sig: "sha256:abc123def456..."
```

## Environment Variables

```bash
# Configure skill proxy (optional)
SKILLPROXY="https://skills.safeclaw.io,direct"

# Private repositories
SKILLPRIVATE="github.com/myorg/*,gitlab.com/internal/*"

# Disable proxy for specific hosts
SKILLNOPROXY="gitee.com,*.corp.example.com"

# Custom directories
SKILL_GLOBAL_DIR="$HOME/.safclaw/skills"
SKILL_WORKSPACE_DIR="$HOME/.safclaw/workspace"
SKILL_CACHE="$HOME/.safeclaw/skill/cache"
```

## Build

```bash
go build -o skill ./cmd/skill
```

## Current Status

✅ Core downloader (GitHub support)
✅ Hook engine with security validation
✅ Installation management
✅ CLI commands (add, remove, list, info, cache, init)
⏳ Full download/install implementation (pending)
⏳ GitLab/Gitee downloaders (pending)
⏳ Proxy support (pending)
⏳ Signature verification (pending)