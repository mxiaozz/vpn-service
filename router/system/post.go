package system

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initSystemPostRouter(pvt wraper.RouterWraper) {
	post := wraper.ExtModule("岗位管理")
	{
		pvt.GET("/post/list", system.PostController.GetPostListPage, post.Ext("system:post:query"))
		pvt.GET("/post/:postId", system.PostController.GetPost, post.Ext("system:post:query"))
		pvt.POST("/post", system.PostController.AddPost, post.Ext("system:post:add"))
		pvt.PUT("/post", system.PostController.UpdatePost, post.Ext("system:post:edit"))
		pvt.DELETE("/post/:postIds", system.PostController.DeletePost, post.Ext("system:post:remove"))
	}
}
