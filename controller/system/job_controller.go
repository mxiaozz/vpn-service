package system

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reugn/go-quartz/quartz"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var JobController = &jobController{}

type jobController struct {
	controller.BaseController
}

func (c *jobController) GetJobListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysJob](ctx)
	if err != nil {
		gb.Logger.Errorln("任务列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var job entity.SysJob
	if err = ctx.ShouldBind(&job); err != nil {
		gb.Logger.Errorln("任务列表获取查询参数失败", err.Error())
		rsp.Fail("任务查询参数格式不正确", ctx)
		return
	}

	// 分页查询
	if err = system.JobService.GetJobListPage(&job, page); err != nil {
		gb.Logger.Errorln("任务列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *jobController) GetJob(ctx *gin.Context) {
	jobId, _ := strconv.ParseInt(ctx.Param("jobId"), 10, 64)
	if jobId == 0 {
		gb.Logger.Errorln("任务列表获取任务详情jobId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if job, err := system.JobService.GetJob(jobId); err != nil {
		gb.Logger.Errorln("任务列表获取岗位详情失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(job, ctx)
	}
}

func (c *jobController) AddJob(ctx *gin.Context) {
	var job entity.SysJob
	if err := ctx.ShouldBind(&job); err != nil {
		gb.Logger.Errorln("增加任务时，获取任务参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := quartz.ValidateCronExpression(job.CronExpression); err != nil {
		gb.Logger.Errorln("增加任务cron不正确", err.Error())
		rsp.Fail("cron表达式格式不正确", ctx)
		return
	}

	// 增补信息
	job.CreateBy = c.GetLoginUser(ctx).User.UserName
	job.CreateTime = time.Now()
	job.Status = "1"

	if err := system.JobService.AddJob(&job); err != nil {
		gb.Logger.Errorln("增加任务失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *jobController) UpdateJob(ctx *gin.Context) {
	var job entity.SysJob
	if err := ctx.ShouldBind(&job); err != nil {
		gb.Logger.Errorln("修改任务时，获取任务参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := quartz.ValidateCronExpression(job.CronExpression); err != nil {
		gb.Logger.Errorln("修改任务cron不正确", err.Error())
		rsp.Fail("cron表达式格式不正确", ctx)
		return
	}

	// 增补信息
	job.UpdateBy = c.GetLoginUser(ctx).User.UserName
	job.UpdateTime = time.Now()

	if err := system.JobService.UpdateJob(&job); err != nil {
		gb.Logger.Errorln("修改任务失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *jobController) ChangeStatus(ctx *gin.Context) {
	var job entity.SysJob
	if err := ctx.ShouldBind(&job); err != nil {
		gb.Logger.Errorln("修改任务状态时，获取任务参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.JobService.ChangeStatus(&job); err != nil {
		gb.Logger.Errorln("修改任务状态失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *jobController) RunJob(ctx *gin.Context) {
	var job entity.SysJob
	if err := ctx.ShouldBind(&job); err != nil {
		gb.Logger.Errorln("执行任务时，获取任务参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	if err := system.JobService.RunJob(&job); err != nil {
		gb.Logger.Errorln("执行任务失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *jobController) DeleteJob(ctx *gin.Context) {
	jobIds := util.NewList(strings.Split(ctx.Param("jobIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(jobIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.JobService.DeleteJobs(jobIds); err != nil {
		gb.Logger.Errorln("删除任务失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
