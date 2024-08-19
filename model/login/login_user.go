package login

import "vpn-web.funcworks.net/model/entity"

type LoginUser struct {
	UserId        int64           `json:"userId"`
	DeptId        int64           `json:"deptId"`
	Token         string          `json:"token"`
	LoginTime     int64           `json:"loginTime"`
	ExpireTime    int64           `json:"expireTime"`
	IpAddress     string          `json:"ipaddr,omitempty"`
	LoginLocation string          `json:"loginLocation,omitempty"`
	Browser       string          `json:"browser,omitempty"`
	Os            string          `json:"os,omitempty"`
	Permissions   map[string]int8 `json:"permissions,omitempty"`
	User          *entity.SysUser `json:"user"`
}
