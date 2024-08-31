package gb

import (
	"context"
)

// 调度任务需实现此接口函数
type WorkFuncion func(params []any, ctx context.Context) (any, error)

// 调度管理器接口
type SchedManager interface {
	// 注册调度任务
	Registry(methodName string, workFunc WorkFuncion) error
	// 调度任务
	ScheduleJob(job SchedJob) error
	// 暂停任务
	PauseJob(job SchedJob) error
	// 恢复任务
	ResumeJob(job SchedJob) error
	// 删除任务
	DeleteJob(job SchedJob) error
	// 执行一次任务
	RunJob(job SchedJob) error
}

// 调度任务
type SchedJob struct {
	JobId          string // 任务ID
	JobName        string // 任务名称
	JobGroup       string // 任务组
	InvokeTarget   string // 调用目标
	CronExpression string // cron表达式
}
