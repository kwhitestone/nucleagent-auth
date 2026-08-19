# nucleagent-auth

认证服务：JWT 签发/验证、用户注册/登录。基于 Prism Fusion 框架。


## Commit Message Language (IRON RULE)

**All commit messages MUST be written in English.** No exceptions.

- Subjects and bodies: English only. No Chinese characters anywhere in the message.
- Type prefixes follow Conventional Commits (`feat:`, `fix:`, `chore:`, `refactor:`, `docs:`, `style:`, `perf:`, `test:`).
- Referencing code identifiers, paths, or domain terms is fine; prose must be English.

Rationale: these repositories are open-sourced on GitHub; Chinese commit messages
make history unreadable to international contributors and pollute git log tooling.

## 构建

> 前置：首次构建需先在 repo 根目录执行 `git submodule update --init` 拉取 prism-fusion，再 `go work sync`。

```bash
cd app/src/server
go work sync
go build ./...
go run main.go        # 启动（需要 MySQL）
```

## 架构约束

- 框架自带 `auth` 插件提供 JWT 登录/注册/刷新
- 自定义认证（OAuth/OpenID）通过额外 addon 实现，不修改框架 auth 插件
- auth provider 通过 `config.yaml` 的 `auth.provider` 切换（builtin / openid）
- JWT secret 与 nucleagent-core 共享（core 本地验证，不远程调用 auth）
- 用户表用框架自带的 User model，不定义业务 model
- CORS：`cors.mode: strict-whitelist`，通过环境变量配置允许的前端来源（WEB_FRONTEND_URL + AUTH_FRONTEND_URL），支持分布式部署

## Addons

| addon | 职责 |
|-------|------|
| auth | 框架内置 JWT 认证（用户名/密码） |
| openid-auth | 可选: OAuth2/OpenID 对接 |
| api-key | API Key 管理 |

## 依赖

- `prism-fusion` (框架) via git submodule + go.work
- MySQL (users, user_roles -- 框架自带表)
- 不依赖 nucleagent-shared（不碰业务表）

## 边界

- **Always**: JWT secret 通过环境变量传入，不硬编码
- **Always**: 密码用 bcrypt
- **Never**: 禁止读写业务表（conversations, agents 等）
- **Never**: 禁止 import nucleagent-core 或 nucleagent-executor
