package router

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/controller/login"
	"vpn-web.funcworks.net/controller/monitor"
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/security"
)

type Router struct {
	*gin.Engine
}

func Init(engine *gin.Engine) {
	initStaticFile(engine)

	sys := extModule("系统登录")
	publicGroup := RouterWraper{engine.Group("/api")}
	publicGroup.GET("/captchaImage", login.Captcha.GetCode, sys.Ext())
	publicGroup.POST("/login", login.UserLogin.Login, sys.Ext())

	privateGroup := RouterWraper{engine.Group("/api")}
	privateGroup.Use(security.JWTAuth())
	privateGroup.POST("/logout", login.UserLogin.Logout, sys.Ext())
	privateGroup.GET("/getInfo", login.UserLogin.GetUserInfo, sys.Ext())
	privateGroup.GET("/getRouters", login.UserLogin.GetRouters, sys.Ext())

	// openvpn
	vpn := extModule("VPN管理")
	privateGroup.GET("/openvpn/getstatus", controller.Openvpn.GetStatus, vpn.Ext())
	privateGroup.GET("/openvpn/getRealtimeStatus", controller.Openvpn.GetRealtimeStatus, vpn.Ext())
	privateGroup.GET("/openvpn/getServerConfig", controller.Openvpn.GetServerConfig, vpn.Ext())
	privateGroup.POST("/openvpn/generateConfig", controller.Openvpn.GenerateConfig, vpn.Ext("openvpn:config:generate"))
	privateGroup.PUT("/openvpn/saveConfig", controller.Openvpn.SaveConfig, vpn.Ext("openvpn:config:save"))
	privateGroup.GET("/openvpn/getPKIStatus", controller.Openvpn.GetPKIStatus, vpn.Ext())
	privateGroup.POST("/openvpn/inikPKI", controller.Openvpn.InikPKI, vpn.Ext("openvpn:pki:init"))
	privateGroup.PUT("/openvpn/resetPKI", controller.Openvpn.ResetPKI, vpn.Ext("openvpn:pki:reset"))
	privateGroup.PUT("/openvpn/optServer", controller.Openvpn.OptServer, vpn.Ext("openvpn:server:opt"))
	privateGroup.GET("/openvpn/getUserCert", controller.Openvpn.GetUserCert, vpn.Ext())
	privateGroup.POST("/openvpn/generateUserCert", controller.Openvpn.GenerateUserCert, vpn.Ext("openvpn:cert:generate"))
	privateGroup.DELETE("/openvpn/revokeUserCert", controller.Openvpn.RevokeUserCert, vpn.Ext("openvpn:cert:revoke"))
	privateGroup.POST("/openvpn/downUserCert", controller.Openvpn.DownUserCert, vpn.Ext("openvpn:cert:download"))

	// 用户管理
	user := extModule("用户管理")
	privateGroup.GET("/system/user/", system.UserController.GetUserInfo, user.Ext("system:user:query"))
	privateGroup.GET("/system/user/:userId", system.UserController.GetUserInfo, user.Ext("system:user:query"))
	privateGroup.POST("/system/user", system.UserController.AddUser, user.Ext("system:user:add"))
	privateGroup.PUT("/system/user", system.UserController.UpdateUser, user.Ext("system:user:edit"))
	privateGroup.DELETE("/system/user/:userIds", system.UserController.DeleteUser, user.Ext("system:user:remove"))
	privateGroup.PUT("/system/user/resetPwd", system.UserController.ResetPassword, user.Ext("system:user:resetPwd"))
	privateGroup.PUT("/system/user/changeStatus", system.UserController.ChangeStatus, user.Ext("system:user:edit"))
	privateGroup.GET("/system/user/authRole/:userId", system.UserController.AuthRole, user.Ext())
	privateGroup.PUT("/system/user/authRole", system.UserController.ChangeUserRoles, user.Ext())
	privateGroup.GET("/system/user/list", system.UserController.GetUserListPage, user.Ext())
	privateGroup.GET("/system/user/deptTree", system.DeptController.GetDeptTree, user.Ext())

	// 菜单管理
	menu := extModule("菜单管理")
	privateGroup.GET("/system/menu/list", system.MenuController.GetMenuList, menu.Ext("system:menu:query"))
	privateGroup.GET("/system/menu/:menuId", system.MenuController.GetMenu, menu.Ext("system:menu:query"))
	privateGroup.POST("/system/menu", system.MenuController.AddMenu, menu.Ext("system:menu:add"))
	privateGroup.PUT("/system/menu", system.MenuController.UpdateMenu, menu.Ext("system:menu:edit"))
	privateGroup.DELETE("/system/menu/:menuId", system.MenuController.DeleteMenu, menu.Ext("system:menu:remove"))
	// 角色新增/编辑时，加载角色对应的菜单，按钮权限
	privateGroup.GET("/system/menu/treeselect", system.MenuController.AddRoleMenuTreeSelect, menu.Ext())
	privateGroup.GET("/system/menu/roleMenuTreeselect/:roleId", system.MenuController.UpdateRoleMenuTreeSelect, menu.Ext())

	// 角色管理
	role := extModule("角色管理")
	privateGroup.GET("/system/role/list", system.RoleController.GetRoleListPage, role.Ext("system:role:query"))
	privateGroup.GET("/system/role/:roleId", system.RoleController.GetRole, role.Ext("system:role:query"))
	privateGroup.POST("/system/role", system.RoleController.AddRole, role.Ext("system:role:add"))
	privateGroup.PUT("/system/role", system.RoleController.UpdateRole, role.Ext("system:role:edit"))
	privateGroup.DELETE("/system/role/:roleIds", system.RoleController.DeleteRole, role.Ext("system:role:remove"))
	privateGroup.PUT("/system/role/changeStatus", system.RoleController.ChangeStatus, role.Ext())
	privateGroup.PUT("/system/role/dataScope", system.RoleController.ChangeRoleDataScope, role.Ext())
	privateGroup.GET("/system/role/deptTree/:roleId", system.RoleController.GetRoleDeptTree, role.Ext())
	privateGroup.GET("/system/role/authUser/allocatedList", system.RoleController.GetRoleUserPage, role.Ext())
	privateGroup.GET("/system/role/authUser/unallocatedList", system.RoleController.GetNotRoleUserPage, role.Ext())
	privateGroup.PUT("/system/role/authUser/selectAll", system.RoleController.AddRoleUsers, role.Ext())
	privateGroup.PUT("/system/role/authUser/cancelAll", system.RoleController.DeleteRoleUsers, role.Ext())
	privateGroup.PUT("/system/role/authUser/cancel", system.RoleController.DeleteRoleUser, role.Ext())

	// 部门管理
	dept := extModule("部门管理")
	privateGroup.GET("/system/dept/list", system.DeptController.GetDeptList, dept.Ext("system:dept:query"))
	privateGroup.GET("/system/dept/:deptId", system.DeptController.GetDept, dept.Ext("system:dept:query"))
	privateGroup.POST("/system/dept", system.DeptController.AddDept, dept.Ext("system:dept:add"))
	privateGroup.GET("/system/dept/list/exclude/:deptId", system.DeptController.GetDeptsExcludeChild, dept.Ext("system:dept:edit"))
	privateGroup.PUT("/system/dept", system.DeptController.UpdateDept, dept.Ext("system:dept:edit"))
	privateGroup.DELETE("/system/dept/:deptId", system.DeptController.DeleteDept, dept.Ext("system:dept:remove"))

	// 岗位管理
	post := extModule("岗位管理")
	privateGroup.GET("/system/post/list", system.PostController.GetPostListPage, post.Ext("system:post:query"))
	privateGroup.GET("/system/post/:postId", system.PostController.GetPost, post.Ext("system:post:query"))
	privateGroup.POST("/system/post", system.PostController.AddPost, post.Ext("system:post:add"))
	privateGroup.PUT("/system/post", system.PostController.UpdatePost, post.Ext("system:post:edit"))
	privateGroup.DELETE("/system/post/:postIds", system.PostController.DeletePost, post.Ext("system:post:remove"))

	// 字典类型管理
	dist := extModule("字典管理")
	privateGroup.GET("/system/dict/type/list", system.DictController.GetDictListPage, dist.Ext("system:dict:query"))
	privateGroup.GET("/system/dict/type/:dictId", system.DictController.GetDict, dist.Ext("system:dict:query"))
	privateGroup.POST("/system/dict/type", system.DictController.AddDict, dist.Ext("system:dict:add"))
	privateGroup.PUT("/system/dict/type", system.DictController.UpdateDict, dist.Ext("system:dict:edit"))
	privateGroup.DELETE("/system/dict/type/:dictIds", system.DictController.DeleteDict, dist.Ext("system:dict:remove"))
	privateGroup.DELETE("/system/dict/type/refreshCache", system.DictController.ReloadDictCache, dist.Ext("system:dict:remove"))
	privateGroup.GET("/system/dict/type/optionselect", system.DictController.GetAllDicts, dist.Ext())

	// 字典数据管理
	privateGroup.GET("/system/dict/data/list", system.DictDataController.GetDictDataListPage, dist.Ext())
	privateGroup.GET("/system/dict/data/type/:dictType", system.DictDataController.GetDictDataByType, dist.Ext())
	privateGroup.GET("/system/dict/data/:dictDataId", system.DictDataController.GetDictData, dist.Ext())
	privateGroup.POST("/system/dict/data", system.DictDataController.AddDictData, dist.Ext("system:dict:add"))
	privateGroup.PUT("/system/dict/data", system.DictDataController.UpdateDictData, dist.Ext("system:dict:add"))
	privateGroup.DELETE("/system/dict/data/:dictDataIds", system.DictDataController.DeleteDictData, dist.Ext("system:dict:add"))

	// 参数管理
	cfg := extModule("参数设置")
	privateGroup.GET("/system/config/list", system.ConfigController.GetConfigListPage, cfg.Ext("system:config:query"))
	privateGroup.GET("/system/config/configKey/:configKey", system.ConfigController.GetConfigByKey, cfg.Ext())
	privateGroup.GET("/system/config/:configId", system.ConfigController.GetConfig, cfg.Ext("system:config:query"))
	privateGroup.POST("/system/config", system.ConfigController.AddConfig, cfg.Ext("system:config:add"))
	privateGroup.PUT("/system/config", system.ConfigController.UpdateConfig, cfg.Ext("system:config:edit"))
	privateGroup.DELETE("/system/config/:configIds", system.ConfigController.DeleteConfig, cfg.Ext("system:config:remove"))
	privateGroup.DELETE("/system/config/refreshCache", system.ConfigController.ReloadConfigCache, cfg.Ext("system:config:remove"))

	// 任务管理
	job := extModule("定时任务")
	privateGroup.GET("/monitor/job/list", system.JobController.GetJobListPage, job.Ext("monitor:job:query"))
	privateGroup.GET("/monitor/job/:jobId", system.JobController.GetJob, job.Ext("monitor:job:query"))
	privateGroup.POST("/monitor/job", system.JobController.AddJob, job.Ext("monitor:job:add"))
	privateGroup.PUT("/monitor/job", system.JobController.UpdateJob, job.Ext("monitor:job:edit"))
	privateGroup.PUT("/monitor/job/changeStatus", system.JobController.ChangeStatus, job.Ext("monitor:job:changeStatus"))
	privateGroup.PUT("/monitor/job/run", system.JobController.RunJob, job.Ext())
	privateGroup.DELETE("/monitor/job/:jobIds", system.JobController.DeleteJob, job.Ext("monitor:job:remove"))

	// 任务日志
	privateGroup.GET("/monitor/jobLog/list", system.JobLogController.GetJobLogListPage, job.Ext())
	privateGroup.GET("/monitor/jobLog/:jobLogId", system.JobLogController.GetJobLog, job.Ext())
	privateGroup.DELETE("/monitor/jobLog/:jobLogIds", system.JobLogController.DeleteJobLog, job.Ext())
	privateGroup.DELETE("/monitor/jobLog/clean", system.JobLogController.CleanJobLogs, job.Ext())

	// 登录日志
	lg := extModule("登录日志")
	privateGroup.GET("/monitor/logininfor/list", monitor.LoginLogController.GetLoginLogListPage, lg.Ext("monitor:logininfor:list"))
	privateGroup.DELETE("/monitor/logininfor/:loginLogIds", monitor.LoginLogController.DeleteLoginLog, lg.Ext("monitor:logininfor:remove"))
	privateGroup.DELETE("/monitor/logininfor/clean", monitor.LoginLogController.CleanLoginLogs, lg.Ext("monitor:logininfor:remove"))
	privateGroup.GET("/monitor/logininfor/unlock/:userName", monitor.LoginLogController.Unlock, lg.Ext("monitor:logininfor:unlock"))

	// 操作日志
	oper := extModule("操作日志")
	privateGroup.GET("/monitor/operlog/list", system.OperLogController.GetOperLogListPage, oper.Ext())
	privateGroup.DELETE("/monitor/operlog/:operLogIds", system.OperLogController.DeleteOperLog, oper.Ext())
	privateGroup.DELETE("/monitor/operlog/clean", system.OperLogController.CleanOperLogs, oper.Ext())

	// 在线用户
	ol := extModule("在线用户")
	privateGroup.GET("/monitor/online/list", monitor.OnlineController.GetOnlineUsers, ol.Ext("monitor:online:list"))
	privateGroup.DELETE("/monitor/online/:userName/:tokenId", monitor.OnlineController.ForceLogout, ol.Ext("monitor:online:forceLogout"))

	// 服务监控
	mt := extModule("服务监控")
	privateGroup.GET("/monitor/server", monitor.ServerController.GetServerInfo, mt.Ext())

	// 缓存列表
	cache := extModule("缓存列表")
	privateGroup.GET("/monitor/cache/getNames", monitor.CacheController.GetCacheNames, cache.Ext("monitor:cache:list"))
	privateGroup.GET("/monitor/cache/getKeys/:cacheName", monitor.CacheController.GetCacheKeys, cache.Ext("monitor:cache:list"))
	privateGroup.GET("/monitor/cache/getValue/:cacheName/:cacheKey", monitor.CacheController.GetCacheValue, cache.Ext("monitor:cache:list"))
	privateGroup.DELETE("/monitor/cache/clearCacheName/:cacheName", monitor.CacheController.ClearCacheName, cache.Ext("monitor:cache:list"))
	privateGroup.DELETE("/monitor/cache/clearCacheKey/:cacheKey", monitor.CacheController.ClearCacheKey, cache.Ext("monitor:cache:list"))
	privateGroup.DELETE("/monitor/cache/clearCacheAll", monitor.CacheController.ClearCacheAll, cache.Ext("monitor:cache:list"))
}

func initStaticFile(engine *gin.Engine) {
	engine.StaticFile("/", "./view/index.html")
	engine.StaticFile("/index.html", "./view/index.html")
	engine.StaticFile("/favicon.ico", "./view/favicon.ico")
	engine.StaticFile("/robots.txt", "./view/robots.txt")
	engine.Static("/html", "./view/html")
	engine.Static("/static", "./view/static")
}
