package apikey

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"whitestone.top/prism-fusion/plugin"
)

// APIKeyPlugin API Key 管理插件。
//
// 复用 auth 插件的路由前缀（/api/v1/addons/auth），使 JWT 中间件作用域自动覆盖。
// BridgeMiddleware 把 JWT 写入的 user_id 桥接到 request context，供 huma handler 读取。
type APIKeyPlugin struct {
	plugin.BasePlugin
}

func init() {
	plugin.Register(&APIKeyPlugin{
		BasePlugin: plugin.BasePlugin{
			PluginName:        "apikey",
			PluginDescription: "API Key 管理 - 创建/列出/删除用户 API Key（哈希存储）",
		},
	})
}

func (p *APIKeyPlugin) Priority() int {
	// 在 auth(10)/rbac(20) 之后执行。
	return 30
}

func (p *APIKeyPlugin) RoutePrefix() string {
	return "/api/v1/addons/auth"
}

func (p *APIKeyPlugin) RegisterRoutes(api huma.API) {
	RegisterRoutes(api)
}

func (p *APIKeyPlugin) Models() []interface{} {
	return []interface{}{&APIKey{}}
}

// Middlewares 作用域中间件：仅对 /api/v1/addons/auth 前缀生效。
// BridgeMiddleware 把 gin context 的 user_id 桥接到 request context。
func (p *APIKeyPlugin) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{BridgeMiddleware()}
}
