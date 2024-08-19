package entity

import (
	"vpn-web.funcworks.net/model"
)

type BaseEntity struct {
	CreateBy   string         `json:"createBy,omitempty"`
	UpdateBy   string         `json:"updateBy,omitempty"`
	CreateTime model.DateTime `json:"createTime"`
	UpdateTime model.DateTime `json:"updateTime"`
	Remark     string         `json:"remark,omitempty"`

	// 前端提交的额外参数
	Params map[string]any `xorm:"-" json:"params,omitempty"`
}
