package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemDeptRouter(pvt wraper.RouterWraper) {
	dept := wraper.ExtModule("部门管理")
	{
		pvt.GET("/dept/list", system.DeptController.GetDeptList, dept.Ext("system:dept:query"))
		pvt.GET("/dept/:deptId", system.DeptController.GetDept, dept.Ext("system:dept:query"))
		pvt.POST("/dept", system.DeptController.AddDept, dept.Ext("system:dept:add"))
		pvt.GET("/dept/list/exclude/:deptId", system.DeptController.GetDeptsExcludeChild, dept.Ext("system:dept:edit"))
		pvt.PUT("/dept", system.DeptController.UpdateDept, dept.Ext("system:dept:edit"))
		pvt.DELETE("/dept/:deptId", system.DeptController.DeleteDept, dept.Ext("system:dept:remove"))
	}
}
