package gb

import (
	"context"

	"vpn-web.funcworks.net/model/entity"
)

// 调度任务需实现此接口函数
type WorkFuncion func(params []any, ctx context.Context) (any, error)

// 调度管理器接口
type SchedManager interface {
	// 注册调度任务
	Registry(methodName string, workFunc WorkFuncion) error
	// 调度任务
	ScheduleJob(jobEntity *entity.SysJob) error
	// 暂停任务
	PauseJob(jobEntity *entity.SysJob) error
	// 恢复任务
	ResumeJob(jobEntity *entity.SysJob) error
	// 删除任务
	DeleteJob(jobEntity *entity.SysJob) error
	// 执行一次任务
	RunJob(jobEntity *entity.SysJob) error
}
