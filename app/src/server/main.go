package main

import (
	"nucleagent_auth/core"
	"nucleagent_auth/global"
	"nucleagent_auth/initialize"

	"go.uber.org/zap"

	_ "whitestone.top/prism-fusion/addons"
	_ "nucleagent_auth/addons"
)

func main() {
	initializeSystem()
	core.RunServer()
}

func initializeSystem() {
	global.PRISM_VP = core.Viper()
	global.PRISM_LOG = core.Zap()
	zap.ReplaceGlobals(global.PRISM_LOG)
	global.PRISM_DB = initialize.Gorm()
	if global.PRISM_DB != nil {
		global.PRISM_LOG.Info("Database connected successfully")
		initialize.InitTables()
	}
	global.PRISM_LOG.Info("nucleagent-auth initialized")
}
