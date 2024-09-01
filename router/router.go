package router

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/router/login"
	"vpn-web.funcworks.net/router/monitor"
	"vpn-web.funcworks.net/router/openvpn"
	"vpn-web.funcworks.net/router/system"
	"vpn-web.funcworks.net/router/wraper"
	"vpn-web.funcworks.net/security"
)

type Router struct {
	*gin.Engine
}

func Init(engine *gin.Engine) {
	pubGroup := wraper.RouterWraper{RouterGroup: engine.Group("/api")}
	pvtGroup := wraper.RouterWraper{RouterGroup: engine.Group("/api")}
	pvtGroup.Use(security.JWTAuth)

	// 登录路由
	login.InitLoginRouter(pubGroup, pvtGroup)

	// 系统管理路由
	system.InitSystemRouter(pvtGroup)

	// 监控管理路由
	monitor.InitMonitorRouter(pvtGroup)

	// VPN管理
	openvpn.InitOpenVpnRouter(pvtGroup)

	initStaticFile(engine)
}

func initStaticFile(engine *gin.Engine) {
	engine.StaticFile("/", "./view/index.html")
	engine.StaticFile("/index.html", "./view/index.html")
	engine.StaticFile("/favicon.ico", "./view/favicon.ico")
	engine.StaticFile("/robots.txt", "./view/robots.txt")
	engine.Static("/html", "./view/html")
	engine.Static("/assets", "./view/assets")
	engine.Static("/static", "./view/static")
}
