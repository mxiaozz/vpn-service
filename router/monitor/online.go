package monitor

import (
	"vpn-web.funcworks.net/controller/monitor"
	"vpn-web.funcworks.net/router/wraper"
)

func initMonitorOnlineRouter(pvt wraper.RouterWraper) {
	ol := wraper.ExtModule("在线用户")
	{
		pvt.GET("/online/list", monitor.OnlineController.GetOnlineUsers, ol.Ext("monitor:online:list"))
		pvt.DELETE("/online/:tokenId", monitor.OnlineController.ForceLogout, ol.Ext("monitor:online:forceLogout"))
	}
}
