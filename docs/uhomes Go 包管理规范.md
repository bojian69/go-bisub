# uhomes Go 包管理规范

**文档版本：** 1.0
**最后更新：** 2025-02-19
**状态：** 草案 (Draft)

---

## 1. 目标与背景

为了规范团队内部 Go 语言项目的依赖管理，提升构建效率，确保代码安全性与一致性，特制定本规范。本规范适用于所有 `git.uhomes.net/uhs-go` 下的 Go 项目。

主要目标：
- **统一依赖源**：确保所有项目使用受控的依赖版本。
- **私有包安全**：规范私有模块的访问与鉴权。
- **版本可追溯**：建立清晰的版本发布与引用机制。
- **环境一致性**：消除本地开发与 CI/CD 环境的差异。

---

## 2. 模块命名规范 (Module Naming)

所有内部 Go 模块必须遵循统一的命名空间：

```go
module git.uhomes.net/uhs-go/<project-name>
```

### 2.1 命名规则
- **域名**：必须以 `git.uhomes.net` 开头。
- **组织名**：统一使用 `uhs-go`。
- **项目名**：使用 `kebab-case`（小写字母 + 短横线），与 Git 仓库名保持一致。
- **多版本**：主版本号 >= v2 时，必须在路径末尾包含版本号（如 `/v2`）。

### 2.2 示例
| 类型 | 仓库地址 | Module Path |
| :--- | :--- | :--- |
| 业务服务 | `git.uhomes.net/uhs-go/user-svc` | `git.uhomes.net/uhs-go/user-svc` |
| 基础库 | `git.uhomes.net/uhs-go/gox` | `git.uhomes.net/uhs-go/gox` |
| Proto 契约 | `git.uhomes.net/uhs-go/proto-contracts` | `git.uhomes.net/uhs-go/proto-contracts` |
| 工具库 (v2) | `git.uhomes.net/uhs-go/utils` | `git.uhomes.net/uhs-go/utils/v2` |

---

## 3. 私有包管理 (Private Modules)

由于我们的代码托管在私有 GitLab 上，Go 工具链默认的代理（如 `proxy.golang.org`）无法访问这些代码。

### 3.1 环境变量配置
所有开发人员和 CI/CD 环境必须配置 `GOPRIVATE` 环境变量。

#### 推荐配置方式（持久化）
建议将配置写入 Shell 配置文件（如 `~/.zshrc` 或 `~/.bash_profile`），以便在所有会话中生效，且不依赖 Go 工具链的全局状态文件。

```bash
# 编辑配置文件
echo 'export GOPRIVATE="git.uhomes.net"' >> ~/.zshrc
echo 'export GOPROXY="https://goproxy.cn,direct"' >> ~/.zshrc

# 使配置生效
source ~/.zshrc
```

#### 备选配置方式（Go 全局设置）
也可以使用 `go env -w` 命令写入 Go 的全局配置文件（通常位于 `os.UserConfigDir`）：

```bash
go env -w GOPRIVATE="git.uhomes.net"
go env -w GOPROXY="https://goproxy.cn,direct"
```

**注意**：`GOPRIVATE` 的作用是告诉 `go` 命令，匹配 `git.uhomes.net` 的模块直接从版本控制系统（Git）获取，不走公共代理（GOPROXY），也不进行 Checksum 校验（GONOSUMDB）。

### 3.2 鉴权配置 (Git Credentials)
Go 命令底层使用 `git` 拉取代码，因此需要配置 Git 的访问权限。

#### 方式一：SSH（推荐，开发环境）
确保本地 SSH Key 已添加到 GitLab 账户。
配置 Git 强制使用 SSH 替代 HTTPS：

```bash
git config --global url."git@git.uhomes.net:".insteadOf "https://git.uhomes.net/"
```

#### 方式二：HTTPS + Token（推荐，CI/CD 环境）
在 CI/CD 流水线中，使用 Access Token：

```bash
git config --global url."https://oauth2:${GITLAB_TOKEN}@git.uhomes.net/".insteadOf "https://git.uhomes.net/"
```

---

## 4. 版本控制策略 (Versioning)

我们严格遵循 [Semantic Versioning 2.0.0](https://semver.org/) 规范。

### 4.1 版本号格式
`v<Major>.<Minor>.<Patch>`

- **Major (主版本)**：不兼容的 API 修改。
- **Minor (次版本)**：向下兼容的功能性新增。
- **Patch (修订号)**：向下兼容的问题修正。

### 4.2 开发阶段版本 (Pseudo-versions)
在未正式发布 Tag 前，可以使用伪版本号引用特定 Commit：

```bash
# 获取最新 commit
go get git.uhomes.net/uhs-go/gox@master

# 获取特定 commit
go get git.uhomes.net/uhs-go/gox@a1b2c3d
```

`go.mod` 中会自动生成类似 `v0.0.0-20250219103000-a1b2c3d4e5f6` 的版本号。

### 4.3 正式发布流程
1. **代码审查**：确保代码通过 Review 且 CI 测试通过。
2. **打标签**：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. **引用更新**：下游项目使用 `go get` 更新依赖。

---

## 5. 依赖管理工作流

### 5.1 初始化/更新依赖
- **添加依赖**：`go get git.uhomes.net/uhs-go/pkg@v1.0.0`
- **整理依赖**：每次提交代码前，**必须**运行：
  ```bash
  go mod tidy
  ```
  这会移除未使用的依赖，并更新 `go.sum`。

### 5.2 依赖升级
- **查看可更新依赖**：
  ```bash
  go list -u -m all
  ```
- **升级所有依赖到最新次版本**：
  ```bash
  go get -u ./...
  go mod tidy
  ```

### 5.3 本地多模块开发 (Workspace Mode)
当需要同时修改 `gox` 基础库和业务服务（如 `user-svc`）时，使用 Go Workspace 避免频繁 push/pull。

1. **创建工作区**：
   ```bash
   mkdir workspace && cd workspace
   go work init
   ```
2. **添加本地模块**：
   ```bash
   go work use ./uhomes-go/gox
   go work use ./uhomes-go/user-svc
   ```
3. **开发**：
   此时 `user-svc` 会直接引用本地的 `gox` 代码，修改即时生效。
4. **提交**：
   - 先提交 `gox` 的更改并 Push。
   - 此时 `user-svc` 的 `go.mod` 仍指向旧版本。
   - 在 `user-svc` 中运行 `go get git.uhomes.net/uhs-go/gox@latest` 更新依赖。
   - 提交 `user-svc`。
   - **注意**：不要提交 `go.work` 文件到仓库（通常加入 `.gitignore`）。

---

## 6. Proto 包管理

Proto 文件的管理遵循 `Uhomes_微服务架构方案_v1.md` 中的定义。

### 6.1 依赖引用
Go 项目不直接引用 `.proto` 文件，而是引用生成的 Go 代码包。

- **公共类型**：`git.uhomes.net/uhs-go/proto-common/...`
- **服务契约**：`git.uhomes.net/uhs-go/proto-contracts/...`

### 6.2 版本锁定
使用 `buf.lock` 锁定 Proto 依赖版本，确保生成代码的一致性。

---

## 7. 最佳实践总结

1. **提交前必做**：运行 `go mod tidy` 确保 `go.mod` 和 `go.sum` 清洁。
2. **不要提交 vendor**：除非有特殊隔离需求，否则不建议提交 `vendor` 目录。
3. **语义化版本**：严格遵守 SemVer，避免破坏性变更导致下游构建失败。
4. **最小依赖原则**：尽量减少引入不必要的第三方库，特别是大型框架。
5. **定期审计**：使用 `govulncheck` 定期检查依赖漏洞。

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```
