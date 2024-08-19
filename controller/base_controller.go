package controller

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/model/login"
)

type BaseController struct {
}

func (c *BaseController) GetLoginUser(ctx *gin.Context) *login.LoginUser {
	user, _ := ctx.Get(cst.SYS_LOGIN_USER_KEY)
	return user.(*login.LoginUser)
}
