---
name: empty-skill-template
description: An empty skill template for creating new skills. Use this template as a starting point when initializing a new skill project with `skill init`.
---

# Empty Skill Template

This is a minimal skill template that provides the basic structure for creating new skills.

## Purpose

When you run `skill init github.com/yourname/yourSkill --template github.com/Safclaw/skill/empty`, this template generates:

- `skill.md` - Skill documentation and usage guide
- `skill.yaml` - Skill configuration and metadata
- `scripts/` - Installation and uninstallation scripts
  - `setup.sh` / `setup.ps1` - Setup scripts for Unix/Windows
  - `unsetup.sh` / `unsetup.ps1` - Cleanup scripts for Unix/Windows

## Structure

```
yourSkill/
├── skill.md          # This file - skill documentation
├── skill.yaml        # Skill configuration
└── scripts/          # Lifecycle hooks
    ├── setup.sh      # Unix setup script
    ├── setup.ps1     # Windows setup script
    ├── unsetup.sh    # Unix cleanup script
    └── unsetup.ps1   # Windows cleanup script
```

## Customization

After generating your skill from this template, you should:

1. **Update `skill.yaml`**:
   - Change `name` to your skill's full name (e.g., `github.com/yourname/yourSkill`)
   - Update `description` with your skill's purpose
   - Add author information
   - Set license type
   - Add relevant tags
   - Configure permissions if needed
   - Add hooks for custom setup/cleanup logic

2. **Customize `skill.md`**:
   - Replace this content with your skill's actual documentation
   - Add usage examples
   - Document any commands or features
   - Include troubleshooting tips

3. **Modify Scripts**:
   - Edit `scripts/setup.sh` and `scripts/setup.ps1` for installation logic
   - Edit `scripts/unsetup.sh` and `scripts/unsetup.ps1` for cleanup logic
   - Add additional scripts if needed

## Example: Creating a Git Helper Skill

```bash
# Initialize new skill from empty template
skill init github.com/yourname/git-helper --template github.com/Safclaw/skill/empty

# Then customize:
# 1. Edit skill.yaml to add git command permissions
# 2. Create skill.md with usage documentation
# 3. Add setup scripts to install git templates
```

## Permissions Template

The empty template starts with no permissions. Add permissions to `skill.yaml` based on your skill's needs:

```yaml
permissions:
  storage:
    - label: "Read project files"
      reason: "Access project source code"
      level: normal
      paths: ["src"]
      rel: workspace
      mod: ["read"]

  network:
    - label: "API Access"
      reason: "Call external APIs"
      level: normal
      domains: ["api.example.com"]

  execution:
    - label: "Run commands"
      reason: "Execute system commands"
      level: normal
      commands: ["git status"]
```

## Hooks Template

Add lifecycle hooks to `skill.yaml`:

```yaml
hooks:
  - stage: "post_add"
    reason: "Initialize skill environment"
    timeout: 60
    scripts:
      - command: "bash"
        platforms: ["linux", "darwin"]
        args: ["./scripts/setup.sh"]
      - command: "powershell"
        platforms: ["windows"]
        args: ["./scripts/setup.ps1"]
```

## Next Steps

After customizing your skill:

1. Test it locally with `skill add ./yourSkill -w`
2. Publish to GitHub or a package repository
3. Share with the community!

## Resources

- [Skill CLI Documentation](https://github.com/Safclaw/skill)
- [Skill YAML Specification](https://github.com/Safclaw/skill/blob/main/ctl/skill.yaml)
- [Example Skills](https://github.com/Safclaw/skills)