# skill
Define skill, and manage skill


## Commands

### 1. Install skill
```bash
# Install to workspace
skill add github.com/Safclaw/skills/read-json

# Install globally
skill add -g github.com/Safclaw/skills/read-json

# Install to specified directory
skill add github.com/Safclaw/skills/read-json --workspace ~/.opencalw/workspace
```

### 2. Uninstall skill
```bash
# Uninstall from workspace
skill remove github.com/Safclaw/skills/read-json

# Uninstall global installation
skill remove -g github.com/Safclaw/skills/read-json

# Uninstall from specified directory
skill remove github.com/Safclaw/skills/read-json --workspace ~/.opencalw/workspace
```

### 3. Initialize a directory as a skill
```bash

# skill init {skillName} [--template github.com/Safclaw/skill/empty]
skill init github.com/xxx/xxSkill
```
After initialization, the directory structure is as follows:
```bash
├── scripts
│   ├── setup.sh # Script executed after installation
│   ├── setup.sp1 # Script executed after installation (Windows)
│   ├── unsetup.sh # Script executed before uninstallation
│   └── unsetup.sp1 # Script executed before uninstallation (Windows)
├── skill.md # Skill entry file
└── skill.yaml # Skill metadata
```