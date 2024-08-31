package entity

import "vpn-web.funcworks.net/model"

type SysUser struct {
	BaseEntity `xorm:"extends"`

	UserId      int64          `xorm:"pk autoincr" json:"userId" form:"userId"`
	DeptId      int64          `json:"deptId,omitempty"`
	UserName    string         `json:"userName,omitempty" form:"userName"`
	NickName    string         `json:"nickName,omitempty"`
	Email       string         `json:"email,omitempty"`
	Phonenumber string         `json:"phonenumber,omitempty" form:"phonenumber"`
	Sex         string         `json:"sex,omitempty"`
	Avatar      string         `json:"avatar,omitempty"`
	Password    string         `json:"password,omitempty"`
	Status      string         `json:"status,omitempty" form:"status"`
	DelFlag     string         `json:"delFlag,omitempty"`
	LoginIp     string         `json:"loginIp,omitempty"`
	LoginDate   model.DateTime `json:"loginDate,omitempty"`
	ValidDay    string         `json:"validDay,omitempty"`

	Dept SysDept `xorm:"-" json:"dept,omitempty"`

	RoleIds []int64 `xorm:"-" json:"roleIds,omitempty" form:"roleIds"`
	PostIds []int64 `xorm:"-" json:"postIds,omitempty" form:"postIds"`
}

func (user *SysUser) IsAdmin() bool {
	return user.UserId == 1
}
