package security

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/model/login"
	"vpn-web.funcworks.net/util/rsp"
)

func PermAuth(handler gin.HandlerFunc, ext ExtInfo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permInfo = ext

		if !authorize(permInfo, ctx) {
			rsp.FailWithCode(cst.HTTP_FORBIDDEN, "无授权", ctx)
			return
		}

		logCtx, err := doBefore(permInfo, ctx)
		if err != nil {
			rsp.Fail(err.Error(), ctx)
			return
		}

		// 执行业务
		handler(ctx)

		doAfter(logCtx, ctx)
	}
}

func authorize(ext ExtInfo, ctx *gin.Context) bool {
	// 没有登记权限，认为无须鉴权
	if len(ext.Perms) == 0 {
		return true
	}

	user, _ := ctx.Get(cst.SYS_LOGIN_USER_KEY)
	loginUser := user.(login.LoginUser)
	userPerms := loginUser.Permissions
	// 用户无任何权限，拒绝
	if len(userPerms) == 0 {
		return false
	}
	// 管理员权限
	if _, ok := userPerms[cst.SYS_ALL_PERMISSION]; ok {
		return true
	}
	// 用户有一项权限，接受
	for _, perm := range ext.Perms {
		if _, ok := userPerms[perm]; ok {
			return true
		}
	}

	// 找不到对应的权限，拒绝
	return false
}
