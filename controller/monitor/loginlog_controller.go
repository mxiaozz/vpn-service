package monitor

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/login"
	"vpn-web.funcworks.net/service/monitor"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var LoginLogController = &loginLogController{}

type loginLogController struct {
	controller.BaseController
}

func (c *loginLogController) GetLoginLogListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysLoginLog](ctx)
	if err != nil {
		gb.Logger.Errorln("登录日志列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var loginLog entity.SysLoginLog
	if err = ctx.ShouldBind(&loginLog); err != nil {
		gb.Logger.Errorln("登录日志列表获取查询参数失败", err.Error())
		rsp.Fail("登录日志查询参数格式不正确", ctx)
		return
	}
	if loginLog.Params == nil {
		loginLog.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			loginLog.Params[k] = v
		}
	}

	// 非管理员只能查看个人日志
	loginUser := c.GetLoginUser(ctx)
	if !loginUser.IsAdmin() {
		loginLog.UserName = loginUser.UserName
	}

	// 分页查询
	if err = monitor.LoginLogService.GetLoginLogListPage(loginLog, page); err != nil {
		gb.Logger.Errorln("登录日志列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *loginLogController) DeleteLoginLog(ctx *gin.Context) {
	loginLogIds := util.NewList(strings.Split(ctx.Param("loginLogIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(loginLogIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := monitor.LoginLogService.DeleteLoginLogs(loginLogIds); err != nil {
		gb.Logger.Errorln("删除登录日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *loginLogController) CleanLoginLogs(ctx *gin.Context) {
	if err := monitor.LoginLogService.CleanLoginLogs(); err != nil {
		gb.Logger.Errorln("清空所有登录日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *loginLogController) Unlock(ctx *gin.Context) {
	userName := ctx.Param("userName")
	if userName == "" {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := login.LoginService.Unlock(userName); err != nil {
		gb.Logger.Errorln("解锁用户失败", userName, err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
