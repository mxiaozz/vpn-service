package system

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var OperLogController = &operLogController{}

type operLogController struct {
	controller.BaseController
}

func (c *operLogController) GetOperLogListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysOperLog](ctx)
	if err != nil {
		gb.Logger.Errorln("操作日志列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var operLog entity.SysOperLog
	if err = ctx.ShouldBind(&operLog); err != nil {
		gb.Logger.Errorln("操作日志列表获取查询参数失败", err.Error())
		rsp.Fail("操作日志查询参数格式不正确", ctx)
		return
	}
	if operLog.Params == nil {
		operLog.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			operLog.Params[k] = v
		}
	}
	if ctx.Query("status") == "" {
		operLog.Status = -1
	}
	if ctx.Query("businessType") == "" {
		operLog.BusinessType = -1
	}

	// 非管理员只能查看个人日志
	loginUser := c.GetLoginUser(ctx)
	if !loginUser.User.IsAdmin() {
		operLog.OperName = loginUser.User.UserName
	}

	// 分页查询
	if err = system.OperLogService.GetOperLogListPage(&operLog, page); err != nil {
		gb.Logger.Errorln("操作日志列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *operLogController) DeleteOperLog(ctx *gin.Context) {
	operLogIds := util.NewList(strings.Split(ctx.Param("operLogIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(operLogIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.OperLogService.DeleteOperLogs(operLogIds); err != nil {
		gb.Logger.Errorln("删除操作日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *operLogController) CleanOperLogs(ctx *gin.Context) {
	if err := system.OperLogService.CleanOperLogs(); err != nil {
		gb.Logger.Errorln("清空所有操作日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
