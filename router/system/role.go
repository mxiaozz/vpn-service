package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemRoleRouter(pvt wraper.RouterWraper) {
	role := wraper.ExtModule("角色管理")
	{
		pvt.GET("/role/list", system.RoleController.GetRoleListPage, role.Ext("system:role:query"))
		pvt.GET("/role/:roleId", system.RoleController.GetRole, role.Ext("system:role:query"))
		pvt.POST("/role", system.RoleController.AddRole, role.Ext("system:role:add"))
		pvt.PUT("/role", system.RoleController.UpdateRole, role.Ext("system:role:edit"))
		pvt.DELETE("/role/:roleIds", system.RoleController.DeleteRole, role.Ext("system:role:remove"))
		pvt.PUT("/role/changeStatus", system.RoleController.ChangeStatus, role.Ext())
		pvt.PUT("/role/dataScope", system.RoleController.ChangeRoleDataScope, role.Ext())
		pvt.GET("/role/deptTree/:roleId", system.RoleController.GetRoleDeptTree, role.Ext())
		pvt.GET("/role/authUser/allocatedList", system.RoleController.GetRoleUserPage, role.Ext())
		pvt.GET("/role/authUser/unallocatedList", system.RoleController.GetNotRoleUserPage, role.Ext())
		pvt.PUT("/role/authUser/selectAll", system.RoleController.AddRoleUsers, role.Ext())
		pvt.PUT("/role/authUser/cancelAll", system.RoleController.DeleteRoleUsers, role.Ext())
		pvt.PUT("/role/authUser/cancel", system.RoleController.DeleteRoleUser, role.Ext())
	}
}
