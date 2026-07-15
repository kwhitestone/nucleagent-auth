# nucleagent-auth

Nucleagent 认证服务。JWT 签发/验证、用户注册/登录、OAuth 对接。

基于 [Prism Fusion](https://github.com/kwhitestone/prism-fusion) 框架构建。

## 结构

```
nucleagent-auth/
├── prism-fusion/              git submodule
├── app/
│   ├── src/
│   │   ├── server/            Go 后端
│   │   │   ├── addons/        业务插件
│   │   │   ├── go.work
│   │   │   ├── go.mod
│   │   │   ├── config.yaml
│   │   │   └── main.go
│   │   └── web/               Vue 前端 (micro-app 子应用)
│   └── Dockerfile
└── README.md
```

## 端口

- 后端: 6670
- 前端: 6678
