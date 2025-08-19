package system

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/openvpn"
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
	user, err := system.UserService.GetSysUserById(userId, true)
	if err != nil {
		gb.Logger.Errorln("获取用户信息失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	user.Password = ""
	data["data"] = user

	// 个人岗位信息
	posts, err = system.PostService.GetUserPostList(userId)
	if err != nil {
		gb.Logger.Errorln("获取个人岗位", err.Error())
		rsp.Fail("获取个人岗位信息失败", ctx)
		return
	}
	data["postIds"] = util.NewList(posts).MapToInt64(func(p entity.SysPost) int64 { return p.PostId })

	// 个人角色信息
	roles, err = system.RoleService.GetUserRoles(userId)
	if err != nil {
		gb.Logger.Errorln("获取用户角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	data["roleIds"] = util.NewList(roles).MapToInt64(func(r entity.SysRole) int64 { return r.RoleId })

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
	user.CreateTime = model.DateTimeNow()
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
	user.UpdateTime = model.DateTimeNow()

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

	for _, id := range userIds {
		if user, err := system.UserService.GetSysUserById(id, false); err != nil {
			gb.Logger.Error(err.Error())
			rsp.Fail(err.Error(), ctx)
		} else {
			if cert, _ := openvpn.OpenvpnService.GetUserCert(user.UserName, false); cert == nil || cert.Name != "" {
				rsp.Fail("请先注销 "+user.UserName+" 用户证书后再删除", ctx)
				return
			}
		}
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

	user, err := system.UserService.GetSysUserById(userId, true)
	if err != nil {
		gb.Logger.Error(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	user.Password = ""

	roles, err := system.RoleService.GetUserRoles(userId)
	if err != nil {
		gb.Logger.Error(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	data := make(map[string]any)
	data["user"] = user
	data["roles"] = roles
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
