package entity

type SysDept struct {
	BaseEntity `xorm:"extends"`

	DeptId    int64  `xorm:"pk autoincr" json:"deptId"`
	ParentId  int64  `json:"parentId"`
	Ancestors string `json:"ancestors"`
	DeptName  string `json:"deptName" form:"deptName"`
	OrderNum  int    `json:"orderNum"`
	Leader    string `json:"leader,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Status    string `json:"status,omitempty" form:"status"`
	DelFlag   string `json:"delFlag,omitempty"`

	Children []SysDept `xorm:"-" json:"children,omitempty"`

	// tree select
	Id    int64  `xorm:"-" json:"id,omitempty"`
	Label string `xorm:"-" json:"label,omitempty"`
}
