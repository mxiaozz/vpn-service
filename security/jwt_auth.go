package security

import (
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/service/login"
	"vpn-web.funcworks.net/util/rsp"
)

func JWTAuth(ctx *gin.Context) {
	loginUser, err := login.TokenService.GetLoginUser(ctx)
	if err != nil {
		if strings.HasSuffix(ctx.Request.URL.Path, "/logout") {
			rsp.Ok(ctx)
		} else {
			rsp.FailWithCode(cst.HTTP_UNAUTHORIZED, "", ctx)
		}
		gb.Logger.Errorf(err.Error())
		ctx.Abort()
		return
	}
	login.TokenService.VerifyToken(loginUser)

	// 将登录用户放置到 request上下文中
	ctx.Set(cst.SYS_LOGIN_USER_KEY, loginUser)

	ctx.Next()
}
