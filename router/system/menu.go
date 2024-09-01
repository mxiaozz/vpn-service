package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemMenuRouter(pvt wraper.RouterWraper) {
	menu := wraper.ExtModule("菜单管理")
	{
		pvt.GET("/menu/list", system.MenuController.GetMenuList, menu.Ext("system:menu:query"))
		pvt.GET("/menu/:menuId", system.MenuController.GetMenu, menu.Ext("system:menu:query"))
		pvt.POST("/menu", system.MenuController.AddMenu, menu.Ext("system:menu:add"))
		pvt.PUT("/menu", system.MenuController.UpdateMenu, menu.Ext("system:menu:edit"))
		pvt.DELETE("/menu/:menuId", system.MenuController.DeleteMenu, menu.Ext("system:menu:remove"))
		// 角色新增/编辑时，加载角色对应的菜单，按钮权限
		pvt.GET("/menu/treeselect", system.MenuController.AddRoleMenuTreeSelect, menu.Ext())
		pvt.GET("/menu/roleMenuTreeselect/:roleId", system.MenuController.UpdateRoleMenuTreeSelect, menu.Ext())
	}
}
