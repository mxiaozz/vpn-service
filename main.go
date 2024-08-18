package main

import (
	_ "vpn-web.funcworks.net/base"
	"vpn-web.funcworks.net/jobtask"
	"vpn-web.funcworks.net/server"
)

func main() {
	// 加载调度任务
	jobtask.LoadSchedJobs()

	// 启动 web 服务
	server.LaunchWeb()
}
