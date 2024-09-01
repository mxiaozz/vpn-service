package monitor

import "vpn-web.funcworks.net/router/wraper"

func InitMonitorRouter(pvt wraper.RouterWraper) {
	pvt = wraper.RouterWraper{RouterGroup: pvt.Group("/monitor")}

	// 任务管理
	initMonitorJobRouter(pvt)

	// 日志管理
	initMonitorLogRouter(pvt)

	// 在线用户
	initMonitorOnlineRouter(pvt)

	// 服务器监控
	initMonitorServerRouter(pvt)

	// 缓存监控
	initMonitorCacheRouter(pvt)
}
