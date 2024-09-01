package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemDictRouter(pvt wraper.RouterWraper) {
	dist := wraper.ExtModule("字典管理")
	// 字典类型管理
	{
		pvt.GET("/dict/type/list", system.DictController.GetDictListPage, dist.Ext("system:dict:query"))
		pvt.GET("/dict/type/:dictId", system.DictController.GetDict, dist.Ext("system:dict:query"))
		pvt.POST("/dict/type", system.DictController.AddDict, dist.Ext("system:dict:add"))
		pvt.PUT("/dict/type", system.DictController.UpdateDict, dist.Ext("system:dict:edit"))
		pvt.DELETE("/dict/type/:dictIds", system.DictController.DeleteDict, dist.Ext("system:dict:remove"))
		pvt.DELETE("/dict/type/refreshCache", system.DictController.ReloadDictCache, dist.Ext("system:dict:remove"))
		pvt.GET("/dict/type/optionselect", system.DictController.GetAllDicts, dist.Ext())
	}
	// 字典数据管理
	{
		pvt.GET("/dict/data/list", system.DictDataController.GetDictDataListPage, dist.Ext())
		pvt.GET("/dict/data/type/:dictType", system.DictDataController.GetDictDataByType, dist.Ext())
		pvt.GET("/dict/data/:dictDataId", system.DictDataController.GetDictData, dist.Ext())
		pvt.POST("/dict/data", system.DictDataController.AddDictData, dist.Ext("system:dict:add"))
		pvt.PUT("/dict/data", system.DictDataController.UpdateDictData, dist.Ext("system:dict:add"))
		pvt.DELETE("/dict/data/:dictDataIds", system.DictDataController.DeleteDictData, dist.Ext("system:dict:add"))
	}
}
