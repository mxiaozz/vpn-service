package monitor

import (
	"vpn-web.funcworks.net/controller/monitor"
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initMonitorLogRouter(pvt wraper.RouterWraper) {
	// 登录日志
	lg := wraper.ExtModule("登录日志")
	{
		pvt.GET("/logininfor/list", monitor.LoginLogController.GetLoginLogListPage, lg.Ext("monitor:logininfor:list"))
		pvt.DELETE("/logininfor/:loginLogIds", monitor.LoginLogController.DeleteLoginLog, lg.Ext("monitor:logininfor:remove"))
		pvt.DELETE("/logininfor/clean", monitor.LoginLogController.CleanLoginLogs, lg.Ext("monitor:logininfor:remove"))
		pvt.GET("/logininfor/unlock/:userName", monitor.LoginLogController.Unlock, lg.Ext("monitor:logininfor:unlock"))
	}
	// 操作日志
	oper := wraper.ExtModule("操作日志")
	{
		pvt.GET("/operlog/list", system.OperLogController.GetOperLogListPage, oper.Ext())
		pvt.DELETE("/operlog/:operLogIds", system.OperLogController.DeleteOperLog, oper.Ext())
		pvt.DELETE("/operlog/clean", system.OperLogController.CleanOperLogs, oper.Ext())
	}
}
