package login

import (
	"vpn-web.funcworks.net/controller/login"
	"vpn-web.funcworks.net/router/wraper"
)

func InitLoginRouter(pub, pvt wraper.RouterWraper) {
	sys := wraper.ExtModule("系统登录")
	{
		pub.GET("/captchaImage", login.Captcha.GetCode, sys.Ext())
		pub.POST("/login", login.UserLogin.Login, sys.Ext())
	}
	{
		pvt.POST("/logout", login.UserLogin.Logout, sys.Ext())
		pvt.GET("/getInfo", login.UserLogin.GetUserInfo, sys.Ext())
		pvt.GET("/getRouters", login.UserLogin.GetRouters, sys.Ext())
	}
}
