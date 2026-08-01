package addons

// 导入此包会触发所有业务插件 init()。
import (
	// API Key 管理插件（创建/列出/删除，哈希存储）
	_ "nucleagent-auth/addons/apikey"
)
