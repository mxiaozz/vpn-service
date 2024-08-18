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

var RoleController = &roleController{}

type roleController struct {
	controller.BaseController
}

// 获取角色列表
func (c *roleController) GetRoleListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysRole](ctx)
	if err != nil {
		gb.Logger.Errorln("角色列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var role entity.SysRole
	if err = ctx.ShouldBind(&role); err != nil {
		gb.Logger.Errorln("角色列表获取查询参数失败", err.Error())
		rsp.Fail("角色查询参数格式不正确", ctx)
		return
	}
	if role.Params == nil {
		role.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			role.Params[k] = v
		}
	}

	// 分页查询
	if err = system.RoleService.GetRoleListPage(&role, page); err != nil {
		gb.Logger.Errorln("角色列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

// 获取角色详情
func (c *roleController) GetRole(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if roleId == 0 {
		gb.Logger.Errorln("角色列表获取角色详情roleId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if role, err := system.RoleService.GetRole(roleId); err != nil {
		gb.Logger.Errorln("角色列表获取角色详情失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(role, ctx)
	}
}

// 添加角色
func (c *roleController) AddRole(ctx *gin.Context) {
	var role entity.SysRole
	if err := ctx.ShouldBind(&role); err != nil {
		gb.Logger.Errorln("增加角色时，获取角色参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	role.CreateBy = c.GetLoginUser(ctx).User.UserName
	role.CreateTime = time.Now()
	role.DelFlag = "0"

	if err := system.RoleService.AddRole(&role); err != nil {
		gb.Logger.Errorln("增加角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 修改角色
func (c *roleController) UpdateRole(ctx *gin.Context) {
	var role entity.SysRole
	if err := ctx.ShouldBind(&role); err != nil {
		gb.Logger.Errorln("修改角色时，获取角色参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	role.UpdateBy = c.GetLoginUser(ctx).User.UserName
	role.UpdateTime = time.Now()

	if err := system.RoleService.UpdateRole(&role); err != nil {
		gb.Logger.Errorln("修改角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 删除角色
func (c *roleController) DeleteRole(ctx *gin.Context) {
	roleIds := util.NewList(strings.Split(ctx.Param("roleIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(roleIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.RoleService.DeleteRoles(roleIds); err != nil {
		gb.Logger.Errorln("删除角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 修改角色状态
func (c *roleController) ChangeStatus(ctx *gin.Context) {
	var role entity.SysRole
	if err := ctx.ShouldBind(&role); err != nil {
		gb.Logger.Errorln("修改角色状态时，获取角色参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.RoleService.ChangeStatus(&role); err != nil {
		gb.Logger.Errorln("修改角色状态失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 获取角色部门树
func (c *roleController) GetRoleDeptTree(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if roleId == 0 {
		gb.Logger.Errorln("获取角色部门树时，获取角色ID参数错误")
		rsp.Fail("参数错误", ctx)
	}

	data := make(map[string]any)
	if depts, err := system.DeptService.GetDeptTree(); err != nil {
		gb.Logger.Errorln("获取角色部门树失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		data["depts"] = depts
	}

	if roleDepts, err := system.DeptService.GetDeptListByRoleId(roleId); err != nil {
		gb.Logger.Errorln("获取角色已关联的部门失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		data["checkedKeys"] = roleDepts
	}

	rsp.Context(ctx).Flat().OkWithData(data)
}

// 修改角色数据权限
func (c *roleController) ChangeRoleDataScope(ctx *gin.Context) {
	var role entity.SysRole
	if err := ctx.ShouldBind(&role); err != nil {
		gb.Logger.Errorln("修改角色数据权限时，获取角色参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.RoleService.ChangeRoleDataScope(&role); err != nil {
		gb.Logger.Errorln("修改角色数据权限失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *roleController) GetRoleUserPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysUser](ctx)
	if err != nil {
		gb.Logger.Errorln("角色已分配用户分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// roleId 参数
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	if roleId == 0 {
		rsp.Fail("参数错误", ctx)
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

	// 分页查询
	if err = system.UserService.GetRoleUserPage(roleId, &user, page); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *roleController) GetNotRoleUserPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysUser](ctx)
	if err != nil {
		gb.Logger.Errorln("角色未分配用户分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// roleId 参数
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	if roleId == 0 {
		rsp.Fail("参数错误", ctx)
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

	// 分页查询
	if err = system.UserService.GetNotRoleUserPage(roleId, &user, page); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *roleController) AddRoleUsers(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	if roleId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	// 解析用户id
	userIds := util.NewList(strings.Split(ctx.Query("userIds"), ",")).
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

	if err := system.RoleService.AddRoleUsers(roleId, userIds); err != nil {
		gb.Logger.Errorln("添加角色用户失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithMessage("添加用户成功", ctx)
	}
}

func (c *roleController) DeleteRoleUsers(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Query("roleId"), 10, 64)
	if roleId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	// 解析用户id
	userIds := util.NewList(strings.Split(ctx.Query("userIds"), ",")).
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

	if err := system.RoleService.DeleteRoleUsers(roleId, userIds); err != nil {
		gb.Logger.Errorln("删除角色用户失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *roleController) DeleteRoleUser(ctx *gin.Context) {
	var data map[string]any
	if err := ctx.ShouldBind(&data); err != nil {
		gb.Logger.Errorln("删除角色用户参数解析失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	roleId, _ := strconv.ParseInt(data["roleId"].(string), 10, 64)
	userId := int64(data["userId"].(float64))

	if err := system.RoleService.DeleteRoleUsers(roleId, []int64{userId}); err != nil {
		gb.Logger.Errorln("删除角色用户失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
