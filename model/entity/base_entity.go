package entity

import "time"

type BaseEntity struct {
	CreateBy   string    `json:"createBy,omitempty"`
	UpdateBy   string    `json:"updateBy,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	Remark     string    `json:"remark,omitempty"`

	// 前端提交的额外参数
	Params map[string]any `xorm:"-" json:"params,omitempty"`
}
