package system

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var PostController = &postController{}

type postController struct {
	controller.BaseController
}

func (c *postController) GetPostListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysPost](ctx)
	if err != nil {
		gb.Logger.Errorln("岗位列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var post entity.SysPost
	if err = ctx.ShouldBind(&post); err != nil {
		gb.Logger.Errorln("岗位列表获取查询参数失败", err.Error())
		rsp.Fail("岗位查询参数格式不正确", ctx)
		return
	}

	// 分页查询
	if err = system.PostService.GetPostListPage(post, page); err != nil {
		gb.Logger.Errorln("岗位列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *postController) GetPost(ctx *gin.Context) {
	postId, _ := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if postId == 0 {
		gb.Logger.Errorln("岗位列表获取岗位详情postId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if post, err := system.PostService.GetPost(postId); err != nil {
		gb.Logger.Errorln("岗位列表获取岗位详情失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(post, ctx)
	}
}

// 添加岗位
func (c *postController) AddPost(ctx *gin.Context) {
	var post entity.SysPost
	if err := ctx.ShouldBind(&post); err != nil {
		gb.Logger.Errorln("增加岗位时，获取岗位参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	post.CreateBy = c.GetLoginUser(ctx).UserName
	post.CreateTime = model.DateTime(time.Now())

	if err := system.PostService.AddPost(post); err != nil {
		gb.Logger.Errorln("增加岗位失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 修改岗位
func (c *postController) UpdatePost(ctx *gin.Context) {
	var post entity.SysPost
	if err := ctx.ShouldBind(&post); err != nil {
		gb.Logger.Errorln("修改岗位时，获取岗位参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	post.UpdateBy = c.GetLoginUser(ctx).UserName
	post.UpdateTime = model.DateTime(time.Now())

	if err := system.PostService.UpdatePost(post); err != nil {
		gb.Logger.Errorln("修改岗位失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 删除角色
func (c *postController) DeletePost(ctx *gin.Context) {
	postIds := util.NewList(strings.Split(ctx.Param("postIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(postIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.PostService.DeletePosts(postIds); err != nil {
		gb.Logger.Errorln("删除角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
