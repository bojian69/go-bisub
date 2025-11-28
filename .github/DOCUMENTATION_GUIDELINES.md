# 📝 文档编写规范

## 文档位置规则

### ✅ 正确做法
- **所有 Markdown 文档统一放在 `docs/` 目录下**
- 项目根目录只保留 `README.md`

### ❌ 错误做法
- ~~在根目录创建其他 `.md` 文件~~
- ~~在其他目录随意创建文档~~

## 文档命名规范

### 文件命名
- 使用英文命名
- 重要文档使用大写：`README.md`, `CHANGELOG.md`, `QUICKSTART.md`
- 指南类文档使用描述性名称：`LOCAL_DEVELOPMENT.md`, `DATABASE_MIGRATION.md`
- 多个单词使用下划线或连字符：`API_REFERENCE.md` 或 `api-reference.md`

### 标题规范
- 使用清晰的层级结构（H1 > H2 > H3）
- H1 标题每个文档只用一次（文档标题）
- 使用 emoji 增强可读性（可选）

## 文档结构

### 基本结构
```markdown
# 文档标题

简短的文档说明（1-2 句话）

## 目录（可选，长文档推荐）

- [章节1](#章节1)
- [章节2](#章节2)

## 章节1

内容...

## 章节2

内容...

## 相关链接

- [其他文档](OTHER_DOC.md)
```

### 必需元素
- 清晰的标题
- 简短的说明
- 代码示例（如果适用）
- 相关文档链接

## 链接规范

### 文档间引用
```markdown
# 在 docs/ 目录内的文档中
[其他文档](OTHER_DOC.md)

# 引用根目录的 README
[项目主页](../README.md)

# 引用 docs 目录（从根目录）
[文档](docs/SOME_DOC.md)
```

### 外部链接
```markdown
[GitHub](https://github.com)
```

## 代码示例规范

### Bash 命令
```markdown
\`\`\`bash
# 注释说明
make dev
\`\`\`
```

### 代码块
```markdown
\`\`\`go
// Go 代码示例
func main() {
    fmt.Println("Hello")
}
\`\`\`
```

### 行内代码
```markdown
使用 `make dev` 启动服务
```

## 文档类型

### 1. 快速入门类
- **START_HERE.md** - 新手入门
- **QUICKSTART.md** - 快速启动
- **COMMANDS.md** - 命令速查

**特点**：
- 简洁明了
- 步骤清晰
- 包含常见问题

### 2. 开发指南类
- **LOCAL_DEVELOPMENT.md** - 本地开发
- **DEPLOYMENT.md** - 部署指南
- **TESTING.md** - 测试指南

**特点**：
- 详细完整
- 包含最佳实践
- 提供故障排查

### 3. 技术文档类
- **API_REFERENCE.md** - API 文档
- **ARCHITECTURE.md** - 架构设计
- **DATABASE_SCHEMA.md** - 数据库设计

**特点**：
- 技术细节
- 设计说明
- 示例代码

### 4. 规范类
- **CODE_STANDARDS.md** - 代码规范
- **GIT_WORKFLOW.md** - Git 工作流
- **NAMING_CONVENTIONS.md** - 命名规范

**特点**：
- 规则明确
- 示例对比
- 工具推荐

## 更新文档

### 添加新文档
1. 在 `docs/` 目录创建文档
2. 更新 `docs/README.md` 索引
3. 在相关文档中添加链接
4. 更新主 `README.md`（如果需要）

### 修改现有文档
1. 保持文档结构一致
2. 更新修改日期（如果有）
3. 检查所有链接是否有效
4. 更新相关引用

## 检查清单

提交文档前请检查：

- [ ] 文档放在 `docs/` 目录下
- [ ] 文件命名符合规范
- [ ] 包含清晰的标题和说明
- [ ] 代码示例格式正确
- [ ] 所有链接有效
- [ ] 更新了 `docs/README.md` 索引
- [ ] 语法和拼写正确
- [ ] 格式统一美观

## 工具推荐

### Markdown 编辑器
- VSCode + Markdown 插件
- Typora
- MacDown (macOS)

### Markdown 检查
- markdownlint
- markdown-link-check

### 预览工具
- GitHub 预览
- VSCode 预览
- Markdown Preview Enhanced

## 示例模板

### 快速入门模板
```markdown
# 文档标题

简短说明

## 前置要求

- 要求1
- 要求2

## 快速开始

\`\`\`bash
# 步骤1
command1

# 步骤2
command2
\`\`\`

## 常见问题

### 问题1
解决方案

## 下一步

- [相关文档](OTHER.md)
```

### 技术文档模板
```markdown
# 技术文档标题

## 概述

简要说明

## 架构设计

详细说明

## 实现细节

### 模块1
说明

### 模块2
说明

## 示例代码

\`\`\`language
代码示例
\`\`\`

## 参考资料

- [链接1](url)
- [链接2](url)
```

## 维护

- 定期检查链接有效性
- 更新过时的内容
- 添加新功能的文档
- 收集用户反馈改进

---

遵循这些规范，让文档更加专业和易用！
