package entity

import "vpn-web.funcworks.net/model"

type SysLoginLog struct {
	InfoId        int64          `xorm:"pk autoincr" json:"infoId"`
	UserName      string         `json:"userName" form:"userName"`
	Status        string         `json:"status" form:"status"`
	Ipaddr        string         `json:"ipaddr" form:"ipaddr"`
	LoginLocation string         `json:"loginLocation"`
	Browser       string         `json:"browser"`
	Os            string         `json:"os"`
	Msg           string         `json:"msg"`
	LoginTime     model.DateTime `json:"loginTime"`

	// 前端提交的额外参数
	Params map[string]any `xorm:"-" json:"params,omitempty"`
}
