package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemUserRouter(pvt wraper.RouterWraper) {
	user := wraper.ExtModule("用户管理")
	{
		// 个人中心
		pvt.GET("/user/profile", system.ProfileController.GetOwnerInfo, user.Ext())
		pvt.PUT("/user/profile", system.ProfileController.UpdateOwnerInfo, user.Ext())
		pvt.PUT("/user/profile/updatePwd", system.ProfileController.UpdateOwnerPassword, user.Ext())
		pvt.POST("/user/profile/avatar", system.ProfileController.UpdateOwnerAvatar, user.Ext())

		// 用户管理
		pvt.GET("/user/", system.UserController.GetUserInfo, user.Ext("system:user:query"))
		pvt.GET("/user/:userId", system.UserController.GetUserInfo, user.Ext("system:user:query"))
		pvt.POST("/user", system.UserController.AddUser, user.Ext("system:user:add"))
		pvt.PUT("/user", system.UserController.UpdateUser, user.Ext("system:user:edit"))
		pvt.DELETE("/user/:userIds", system.UserController.DeleteUser, user.Ext("system:user:remove"))
		pvt.PUT("/user/resetPwd", system.UserController.ResetPassword, user.Ext("system:user:resetPwd"))
		pvt.PUT("/user/changeStatus", system.UserController.ChangeStatus, user.Ext("system:user:edit"))
		pvt.GET("/user/authRole/:userId", system.UserController.AuthRole, user.Ext())
		pvt.PUT("/user/authRole", system.UserController.ChangeUserRoles, user.Ext())
		pvt.GET("/user/list", system.UserController.GetUserListPage, user.Ext())
		pvt.GET("/user/deptTree", system.DeptController.GetDeptTree, user.Ext())
	}
}
