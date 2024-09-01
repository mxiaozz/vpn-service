package monitor

import (
	"vpn-web.funcworks.net/controller/system"
	"vpn-web.funcworks.net/router/wraper"
)

func initMonitorJobRouter(pvt wraper.RouterWraper) {
	job := wraper.ExtModule("定时任务")
	// 任务管理
	{
		pvt.GET("/job/list", system.JobController.GetJobListPage, job.Ext("monitor:job:query"))
		pvt.GET("/job/:jobId", system.JobController.GetJob, job.Ext("monitor:job:query"))
		pvt.POST("/job", system.JobController.AddJob, job.Ext("monitor:job:add"))
		pvt.PUT("/job", system.JobController.UpdateJob, job.Ext("monitor:job:edit"))
		pvt.PUT("/job/changeStatus", system.JobController.ChangeStatus, job.Ext("monitor:job:changeStatus"))
		pvt.PUT("/job/run", system.JobController.RunJob, job.Ext())
		pvt.DELETE("/job/:jobIds", system.JobController.DeleteJob, job.Ext("monitor:job:remove"))
	}
	// 任务日志
	{
		pvt.GET("/jobLog/list", system.JobLogController.GetJobLogListPage, job.Ext())
		pvt.GET("/jobLog/:jobLogId", system.JobLogController.GetJobLog, job.Ext())
		pvt.DELETE("/jobLog/:jobLogIds", system.JobLogController.DeleteJobLog, job.Ext())
		pvt.DELETE("/jobLog/clean", system.JobLogController.CleanJobLogs, job.Ext())
	}
}
