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

var JobLogController = &jobLogController{}

type jobLogController struct {
	controller.BaseController
}

func (c *jobLogController) GetJobLogListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysJobLog](ctx)
	if err != nil {
		gb.Logger.Errorln("任务日志列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var jobLog entity.SysJobLog
	if err = ctx.ShouldBind(&jobLog); err != nil {
		gb.Logger.Errorln("任务日志列表获取查询参数失败", err.Error())
		rsp.Fail("任务日志查询参数格式不正确", ctx)
		return
	}
	if jobLog.Params == nil {
		jobLog.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			jobLog.Params[k] = v
		}
	}

	// 分页查询
	if err = system.JobLogService.GetJobLogListPage(jobLog, page); err != nil {
		gb.Logger.Errorln("任务日志列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *jobLogController) GetJobLog(ctx *gin.Context) {
	jobLogId, _ := strconv.ParseInt(ctx.Param("jobLogId"), 10, 64)
	if jobLogId == 0 {
		gb.Logger.Errorln("任务日志列表获取任务详情jobLogId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if jobLog, err := system.JobLogService.GetJobLog(jobLogId); err != nil {
		gb.Logger.Errorln("任务日志列表获取日志详情失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(jobLog, ctx)
	}
}

func (c *jobLogController) DeleteJobLog(ctx *gin.Context) {
	jobLogIds := util.NewList(strings.Split(ctx.Param("jobLogIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(jobLogIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.JobLogService.DeleteJobLogs(jobLogIds); err != nil {
		gb.Logger.Errorln("删除任务日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *jobLogController) CleanJobLogs(ctx *gin.Context) {
	if err := system.JobLogService.CleanJobLogs(); err != nil {
		gb.Logger.Errorln("清空所有任务日志失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
