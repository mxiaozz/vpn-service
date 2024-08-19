package system

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var MenuController = &menuController{}

type menuController struct {
	controller.BaseController
}

// 菜单管理列表
func (c *menuController) GetMenuList(ctx *gin.Context) {
	var menu entity.SysMenu
	if err := ctx.ShouldBindQuery(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	if !util.IsAdminId(loginUser.UserId) {
		menu.Params["userId"] = loginUser.UserId
	}

	if menus, err := system.MenuService.GetMenuList(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(menus, ctx)
	}
}

// 菜单管理详情
func (c *menuController) GetMenu(ctx *gin.Context) {
	menuId, _ := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if menuId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if menu, err := system.MenuService.GetMenu(menuId); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(menu, ctx)
	}
}

// 菜单管理新增
func (c *menuController) AddMenu(ctx *gin.Context) {
	var menu entity.SysMenu
	if err := ctx.ShouldBind(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if cst.MENU_YES_FRAME == menu.IsFrame && !util.IsHttp(menu.Path) {
		rsp.Fail("新增菜单'"+menu.MenuName+"'失败，地址必须以http(s)://开头", ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	menu.CreateBy = loginUser.User.UserName
	menu.CreateTime = model.DateTime(time.Now())

	if err := system.MenuService.AddMenu(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 菜单管理修改
func (c *menuController) UpdateMenu(ctx *gin.Context) {
	var menu entity.SysMenu
	if err := ctx.ShouldBind(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if cst.MENU_YES_FRAME == menu.IsFrame && !util.IsHttp(menu.Path) {
		rsp.Fail("修改菜单'"+menu.MenuName+"'失败，地址必须以http(s)://开头", ctx)
		return
	}
	if menu.MenuId == menu.ParentId {
		rsp.Fail("修改菜单'"+menu.MenuName+"'失败，上级菜单不能选择自己", ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	menu.UpdateBy = loginUser.User.UserName
	menu.UpdateTime = model.DateTime(time.Now())

	if err := system.MenuService.UpdateMenu(&menu); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 菜单管理删除
func (c *menuController) DeleteMenu(ctx *gin.Context) {
	menuId, _ := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if menuId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.MenuService.DeleteMenu(menuId); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

// 角色管理，角色新增时菜单选择框列表
func (c *menuController) AddRoleMenuTreeSelect(ctx *gin.Context) {
	if menus, err := system.MenuService.GetRolerMenuTreeSelect(c.GetLoginUser(ctx)); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		rsp.OkWithData(menus, ctx)
	}
}

// 角色管理，角色编辑菜单选择框列表
func (c *menuController) UpdateRoleMenuTreeSelect(ctx *gin.Context) {
	roleId, _ := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if roleId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	data := make(map[string]interface{})
	if menus, err := system.MenuService.GetRolerMenuTreeSelect(c.GetLoginUser(ctx)); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		data["menus"] = menus
	}
	if menuIds, err := system.MenuService.GetMenuListByRoleId(roleId); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		data["checkedKeys"] = menuIds
	}

	rsp.Context(ctx).Flat().OkWithData(data)
}
