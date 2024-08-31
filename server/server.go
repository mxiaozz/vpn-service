package server

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/router"
	"vpn-web.funcworks.net/util"
)

func LaunchWeb() {
	// 设置gin模式
	gin.SetMode(util.If(gb.Config.Server.DevMode, gin.DebugMode, gin.ReleaseMode))

	// web engine
	engine := gin.New()
	engine.Use(gin.LoggerWithFormatter(ginLogFormat), gin.Recovery())

	// 注册路由
	router.Init(engine)

	// 404
	noRouteHandle(engine)

	// 启动服务
	engine.Run(getAddr())
}

func noRouteHandle(engine *gin.Engine) {
	// 前后端分离，F5刷新时路径为前端路由器，后台将404，将调转至首页
	engine.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" {
			ctx.File("./view/index.html")
		} else {
			ctx.JSON(404, gin.H{
				"code":    404,
				"message": "Page Not Found",
			})
		}
	})
}

func getAddr() string {
	host := gb.Config.Server.Host
	if host == "" && gb.Config.Server.DevMode {
		host = "127.0.0.1"
	}
	return host + ":" + gb.Config.Server.Port
}
