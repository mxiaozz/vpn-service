package entity

type SysRole struct {
	BaseEntity `xorm:"extends"`

	RoleId            int64  `xorm:"pk autoincr" json:"roleId"`
	RoleName          string `json:"roleName,omitempty" form:"roleName"`
	RoleKey           string `json:"roleKey,omitempty" form:"roleKey"`
	RoleSort          int    `json:"roleSort,omitempty"`
	DataScope         string `json:"dataScope,omitempty"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly,omitempty"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly,omitempty"`
	Status            string `json:"status,omitempty" form:"status"`
	DelFlag           string `json:"delFlag,omitempty"`

	MenuIds []int64 `xorm:"-" json:"menuIds,omitempty"`
	DeptIds []int64 `xorm:"-" json:"deptIds,omitempty"`
}

func (role *SysRole) IsAdmin() bool {
	return role.RoleId == 1
}
