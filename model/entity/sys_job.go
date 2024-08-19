package entity

import "vpn-web.funcworks.net/model"

type SysJob struct {
	BaseEntity `xorm:"extends"`

	JobId          int64  `xorm:"pk autoincr" json:"jobId" form:"jobId"` // 任务ID
	JobName        string `json:"jobName" form:"jobName"`                // 任务名称
	JobGroup       string `json:"jobGroup" form:"jobGroup"`              // 任务组名
	InvokeTarget   string `json:"invokeTarget" form:"invokeTarget"`      // 调用目标字符串
	CronExpression string `json:"cronExpression" form:"cronExpression"`  // cron执行表达式
	MisfirePolicy  string `json:"misfirePolicy" form:"misfirePolicy"`    // cron计划策略
	Concurrent     string `json:"concurrent" form:"concurrent"`          // 是否并发执行（0允许 1禁止）
	Status         string `json:"status" form:"status"`                  // 任务状态（0正常 1暂停）

	NextValidTime model.DateTime `xorm:"-" json:"nextValidTime"`
}
