package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util/rsp"
)

var DeptController = &deptController{}

type deptController struct {
	controller.BaseController
}

func (c *deptController) GetDeptTree(ctx *gin.Context) {
	if depts, err := system.DeptService.GetDeptTree(); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(depts, ctx)
	}
}

func (c *deptController) GetDeptList(ctx *gin.Context) {
	var dept entity.SysDept
	if err := ctx.ShouldBind(&dept); err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	if depts, err := system.DeptService.GetDeptList(dept); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(depts, ctx)
	}
}

func (c *deptController) GetDept(ctx *gin.Context) {
	deptId, _ := strconv.ParseInt(ctx.Param("deptId"), 10, 64)
	if deptId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if dept, err := system.DeptService.GetDept(deptId); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(dept, ctx)
	}
}

func (c *deptController) AddDept(ctx *gin.Context) {
	var dept entity.SysDept
	if err := ctx.ShouldBind(&dept); err != nil {
		gb.Logger.Errorln("增加部门时，获取部门参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dept.CreateBy = c.GetLoginUser(ctx).UserName
	dept.CreateTime = model.DateTimeNow()
	dept.DelFlag = "0"

	if err := system.DeptService.AddDept(dept); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *deptController) GetDeptsExcludeChild(ctx *gin.Context) {
	deptId, _ := strconv.ParseInt(ctx.Param("deptId"), 10, 64)
	if deptId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if depts, err := system.DeptService.GetDeptsExcludeChild(deptId); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(depts, ctx)
	}
}

func (c *deptController) UpdateDept(ctx *gin.Context) {
	var dept entity.SysDept
	if err := ctx.ShouldBind(&dept); err != nil {
		gb.Logger.Errorln("编辑部门时，获取部门参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dept.UpdateBy = c.GetLoginUser(ctx).UserName
	dept.UpdateTime = model.DateTimeNow()

	if err := system.DeptService.UpdateDept(dept); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *deptController) DeleteDept(ctx *gin.Context) {
	deptId, _ := strconv.ParseInt(ctx.Param("deptId"), 10, 64)
	if deptId == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.DeptService.DeleteDept(deptId); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
