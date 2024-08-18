package entity

import "time"

type SysJobLog struct {
	JobLogId      int64     `xorm:"pk autoincr" json:"jobLogId"`
	JobName       string    `json:"jobName" form:"jobName"`
	JobGroup      string    `json:"jobGroup" form:"jobGroup"`
	InvokeTarget  string    `json:"invokeTarget"`
	JobMessage    string    `json:"jobMessage"`
	Status        string    `json:"status" form:"status"` // 执行状态（0正常 1失败）
	ExceptionInfo string    `json:"exceptionInfo"`
	CreateTime    time.Time `json:"createTime"`

	// 前端提交的额外参数
	Params map[string]any `xorm:"-" json:"params,omitempty"`
}
