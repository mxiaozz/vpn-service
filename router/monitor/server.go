package monitor

import (
	"vpn-web.funcworks.net/controller/monitor"
	"vpn-web.funcworks.net/router/wraper"
)

func initMonitorServerRouter(pvt wraper.RouterWraper) {
	mt := wraper.ExtModule("服务监控")
	{
		pvt.GET("/server", monitor.ServerController.GetServerInfo, mt.Ext())
	}
}
