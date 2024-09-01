package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemConfigRouter(pvt wraper.RouterWraper) {
	cfg := wraper.ExtModule("参数设置")
	{
		pvt.GET("/config/list", system.ConfigController.GetConfigListPage, cfg.Ext("system:config:query"))
		pvt.GET("/config/configKey/:configKey", system.ConfigController.GetConfigByKey, cfg.Ext())
		pvt.GET("/config/:configId", system.ConfigController.GetConfig, cfg.Ext("system:config:query"))
		pvt.POST("/config", system.ConfigController.AddConfig, cfg.Ext("system:config:add"))
		pvt.PUT("/config", system.ConfigController.UpdateConfig, cfg.Ext("system:config:edit"))
		pvt.DELETE("/config/:configIds", system.ConfigController.DeleteConfig, cfg.Ext("system:config:remove"))
		pvt.DELETE("/config/refreshCache", system.ConfigController.ReloadConfigCache, cfg.Ext("system:config:remove"))
	}
}
