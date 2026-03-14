---
name: skill-manager
description: SafeClaw skill CLI 工具使用指南。当用户提到 skill add/remove/list/init/info/cache 命令，或需要安装、卸载、查看、初始化 skill 包时触发。
---

包名格式：`{host}/{namespace}/{name}[@{version}]`

## 命令

```bash
# 安装（-w 工作区 | -g 全局 | --workspace-path <dir>）
skill add github.com/Safclaw/skills/read-json -w
skill add -g github.com/Safclaw/skills/read-json@v1.0.0

# 卸载（参数同上）
skill remove github.com/Safclaw/skills/read-json -w

# 列出 / 查看
skill list -w / -g / --workspace-path <dir>
skill info github.com/Safclaw/skills/read-json

# 初始化新 skill（生成 skill.md, skill.yaml, scripts/）
skill init github.com/yourname/yourSkill [--template github.com/Safclaw/skill/empty]

# 缓存
skill cache clean [--all] | verify
```

## 环境变量

| 变量 | 用途 |
|------|------|
| `SKILLPROXY` | 代理，如 `https://skills.safeclaw.io,direct` |
| `SKILLPRIVATE` | 私有仓库规则，如 `github.com/myorg/*` |
| `SKILLNOPROXY` | 跳过代理域名 |
| `SKILL_GLOBAL_DIR` | 全局安装目录 |
| `SKILL_WORKSPACE_DIR` | 工作区目录 |
| `SKILL_CACHE` | 缓存目录 |

## 故障排查

- 下载失败 → 检查 `SKILLPROXY`
- 校验失败 → `skill cache verify`，必要时 `skill cache clean --all` 后重装
- GitLab/Gitee、代理、签名验证功能仍在开发中