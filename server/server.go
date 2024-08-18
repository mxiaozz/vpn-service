package server

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/router"
)

func LaunchWeb() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.LoggerWithFormatter(ginLogFormat), gin.Recovery())
	router.Init(engine)

	host := gb.Config.Server.Host
	if host == "" && gb.Config.Server.DevMode {
		gin.SetMode(gin.DebugMode)
		host = "127.0.0.1"
	}
	port := gb.Config.Server.Port
	addr := host + ":" + port
	engine.Run(addr)
}
