package jobtask

import "vpn-web.funcworks.net/gb"

func registryJob() {
	// 定时删除job日志
	gb.Sched.Registry("cleanJobLog", cleanJobLog)

	// 定时更新用户证书有效期剩余天数
	gb.Sched.Registry("refreshUserCertStatus", refreshUserCertStatus)

	// 定时刷新 OpenVPN 服务状态
	gb.Sched.Registry("refreshVpnStatus", vpnStausSchedule.refreshVpnStatus)
}
