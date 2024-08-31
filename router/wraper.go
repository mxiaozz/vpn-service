package router

import (
	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/security"
)

type RouterWraper struct {
	*gin.RouterGroup
}

func (wrap *RouterWraper) GET(relativePath string, handlers gin.HandlerFunc, ext security.ExtInfo) gin.IRoutes {
	return wrap.RouterGroup.GET(relativePath, security.PermAuth(handlers, ext))
}

func (wrap *RouterWraper) POST(relativePath string, handlers gin.HandlerFunc, ext security.ExtInfo) gin.IRoutes {
	return wrap.RouterGroup.POST(relativePath, security.PermAuth(handlers, ext))
}

func (wrap *RouterWraper) PUT(relativePath string, handlers gin.HandlerFunc, ext security.ExtInfo) gin.IRoutes {
	return wrap.RouterGroup.PUT(relativePath, security.PermAuth(handlers, ext))
}

func (wrap *RouterWraper) DELETE(relativePath string, handlers gin.HandlerFunc, ext security.ExtInfo) gin.IRoutes {
	return wrap.RouterGroup.DELETE(relativePath, security.PermAuth(handlers, ext))
}

// 路由扩展信息
// module 应用于日志记录
// perms  为用户授权权限标识
func extModule(module string) security.ExtInfo {
	return security.ExtInfo{Module: module}
}
