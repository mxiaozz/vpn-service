package base

import "vpn-web.funcworks.net/gb"

func initDefaultConfig() {
	gb.Config.Server.Port = "8080"

	gb.Config.Login.MaxRetryCount = 5
	gb.Config.Login.LockTime = 10
}
