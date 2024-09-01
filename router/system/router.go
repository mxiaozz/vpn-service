package system

import "vpn-web.funcworks.net/router/wraper"

func InitSystemRouter(pvt wraper.RouterWraper) {
	pvt = wraper.RouterWraper{RouterGroup: pvt.Group("/system")}

	// 用户管理/个人中心
	initSystemUserRouter(pvt)

	// 菜单管理
	initSystemMenuRouter(pvt)

	// 角色管理
	initSystemRoleRouter(pvt)

	// 部门管理
	initSystemDeptRouter(pvt)

	// 岗位管理
	initSystemPostRouter(pvt)

	// 字典管理
	initSystemDictRouter(pvt)

	// 参数管理
	initSystemConfigRouter(pvt)
}
