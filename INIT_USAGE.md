# skill init 命令增强说明

## 功能概述

`skill init` 命令已增强，支持从模板创建 skill 项目。主要特性：

1. **在当前目录初始化** - 不再创建子目录，而是在当前工作目录生成文件
2. **默认使用 empty 模板** - empty 模板已编译到可执行程序中，无需外部依赖
3. **支持自定义模板** - 通过 `--template` 参数指定本地或远程模板
4. **模板验证** - 自动检查模板是否符合 skill 规范（必须包含 skill.yaml）
5. **自动更新模块名** - 将模板中的 name 字段替换为命令行指定的模块名

## 使用方法

### 1. 使用默认模板

```bash
cd /path/to/new-skill
skill init github.com/myorg/my-skill
```

这会在当前目录创建 skill 项目，使用默认的空模板（已嵌入到可执行程序中）。

### 2. 使用本地模板

```bash
# 绝对路径
skill init github.com/myorg/my-skill --template /path/to/template

# 相对路径
skill init github.com/myorg/my-skill --template ./my-template
skill init github.com/myorg/my-skill --template ../templates/empty
```

### 3. 使用远程仓库模板

```bash
# 从 GitHub 仓库的特定目录
skill init github.com/myorg/my-skill --template github.com/Safclaw/skill/empty@main
skill init github.com/myorg/my-skill --template github.com/user/repo/template@v1.0.0
```

格式：`github.com/{owner}/{repo}/{subdir}@{version}`

- `{owner}/{repo}`: GitHub 仓库
- `{subdir}`: 仓库中的子目录（可选）
- `{version}`: 分支名、tag 或 commit hash（可选，默认为 main）

## 模板规范

一个有效的 skill 模板必须满足：

1. 包含 `skill.yaml` 文件
2. `skill.yaml` 中必须包含 `name:` 字段
3. 可以包含以下可选文件和目录：
   - `skill.md` - skill 实现文件
   - `scripts/` - 安装和卸载脚本目录
   - 其他 skill 所需的文件

## 示例模板结构

```
my-template/
├── skill.yaml      # 必需，包含 skill 元数据
├── skill.md        # 可选，skill 实现
├── README.md       # 可选，项目说明
└── scripts/        # 可选，脚本目录
    ├── setup.sh
    ├── unsetup.sh
    ├── setup.ps1
    └── unsetup.ps1
```

## 实现细节

### 核心组件

1. **TemplateManager** (`pkg/template/manager.go`)
   - 模板下载和管理
   - 本地/远程模板识别
   - 嵌入模板支持（使用 Go embed）
   - 模板验证
   - 文件复制和模块名更新

2. **Templates Package** (`pkg/templates/embed.go`)
   - 使用 `//go:embed` 指令将 empty 模板编译到二进制文件中
   - 提供嵌入文件系统访问接口

3. **init.go** (`cmd/skill/init.go`)
   - 简化的初始化逻辑
   - 调用 TemplateManager 复制模板
   - 用户交互和提示

### 工作流程

1. 解析命令行参数
2. 确定模板路径（默认使用嵌入的 empty 模板，或 `--template` 指定）
3. 判断模板类型（嵌入、本地或远程）
4. 从嵌入文件系统/本地/网络下载模板文件到当前目录
5. 验证模板有效性
6. 更新 `skill.yaml` 中的 `name` 字段
7. 输出成功消息和后续步骤提示

## 测试

```bash
# 测试默认模板
mkdir test1 && cd test1
skill init github.com/test/my-skill

# 测试本地模板
mkdir test2 && cd test2
skill init github.com/test/my-skill --template /path/to/empty

# 测试远程模板
mkdir test3 && cd test3
skill init github.com/test/my-skill --template github.com/Safclaw/skill/empty@main
```

检查生成的文件：
- `skill.yaml` - name 字段应更新为命令行指定的模块名
- `skill.md` - 空文件或模板内容
- `scripts/` - 安装和卸载脚本

## 注意事项

1. **权限问题** - 确保对模板目录有读取权限
2. **网络访问** - 使用远程模板时需要访问 GitHub
3. **版本指定** - 远程模板建议使用明确的版本号或 tag，避免使用 latest
4. **模板验证** - 无效的模板会被拒绝，并显示错误信息

## Bug 修复

本次更新修复了以下问题：

- ✅ skill init 在当前目录初始化，而不是创建子目录
- ✅ 第一个参数是模块名，不是目录路径
- ✅ 支持从 empty 目录作为模板创建
- ✅ 支持通过 --template 指定本地或远程模板
- ✅ 自动验证模板是否符合 skill 规范
