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

var UserController = &userController{}

type userController struct {
	controller.BaseController
}

// 获取用户列表
func (c *userController) GetUserListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysUser](ctx)
	if err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 获取查询参数
	var user entity.SysUser
	err = ctx.ShouldBind(&user)
	if err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	if user.Params == nil {
		user.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			user.Params[k] = v
		}
	}

	// 分页查询
	if err = system.UserService.GetUserListPage(user, page); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

// 根据用户编号获取详细信息
func (c *userController) GetUserInfo(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.Param("userId"), 10, 64)

	data := make(map[string]any)

	// 获取角色
	roles, err := system.RoleService.GetAllRoles()
	if err != nil {
		gb.Logger.Errorln("获取角色", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	if util.IsAdminId(userId) {
		data["roles"] = roles
	} else {
		data["roles"] = util.NewList(roles).Filter(func(r entity.SysRole) bool { return !r.IsAdmin() })
	}

	// 获取岗位
	posts, err := system.PostService.GetALlPosts()
	if err != nil {
		gb.Logger.Errorln("获取岗位", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	data["posts"] = posts

	if userId == 0 {
		rsp.Context(ctx).Flat().OkWithData(data)
		return
	}

	// 个人信息详情
	sysUser, err := system.UserService.GetSysUserById(userId, true)
	if err != nil {
		gb.Logger.Errorln("获取用户信息失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	data["data"] = sysUser

	// 个人岗位信息
	posts, err = system.PostService.GetUserPostList(userId)
	if err != nil {
		gb.Logger.Errorln("获取个人岗位", err.Error())
		rsp.Fail("获取个人岗位信息失败", ctx)
		return
	}
	data["postIds"] = util.NewList(posts).MapToInt64(func(p entity.SysPost) int64 { return p.PostId })

	// 个人角色信息
	data["roleIds"] = util.NewList(sysUser.Roles).MapToInt64(func(r entity.SysRole) int64 { return r.RoleId })

	rsp.Context(ctx).Flat().OkWithData(data)
}

func (c *userController) AddUser(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBind(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	user.CreateBy = c.GetLoginUser(ctx).UserName
	user.CreateTime = model.DateTime(time.Now())
	user.DelFlag = "0"

	// 数据库增加
	if err := system.UserService.AddUser(user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBind(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	user.UpdateBy = c.GetLoginUser(ctx).UserName
	user.UpdateTime = model.DateTime(time.Now())

	if err := system.UserService.UpdateUser(user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	// 解析用户id
	userIds := util.NewList(strings.Split(ctx.Param("userIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(userIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.UserService.DeleteUser(userIds); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.Ok(ctx)
}

func (c *userController) ResetPassword(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBind(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.UserService.ResetPassword(user); err != nil {
		gb.Logger.Errorf("修改密码失败：%d, %s", user.UserId, err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *userController) ChangeStatus(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBind(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.UserService.ChangeStatus(user); err != nil {
		gb.Logger.Errorf("修改状态失败：%d, %s", user.UserId, err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		rsp.Ok(ctx)
	}
}

func (c *userController) AuthRole(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if userId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	sysUser, err := system.UserService.GetSysUserById(userId, true)
	if err != nil {
		gb.Logger.Error(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	data := make(map[string]any)
	data["user"] = sysUser
	data["roles"] = sysUser.Roles
	sysUser.Roles = nil
	rsp.Context(ctx).Flat().OkWithData(data)
}

func (c *userController) ChangeUserRoles(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBindQuery(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	user.RoleIds = util.NewList(user.RoleIds).Filter(func(id int64) bool { return id > 0 })
	if len(user.RoleIds) == 0 {
		rsp.Ok(ctx)
		return
	}

	if err := system.UserService.ChangeUserRoles(user.UserId, user.RoleIds); err != nil {
		gb.Logger.Errorf("修改用户角色失败：%d, %s", user.UserId, err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
