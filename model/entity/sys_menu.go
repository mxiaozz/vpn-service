package entity

type SysMenu struct {
	BaseEntity `xorm:"extends"`

	MenuId    int64  `xorm:"pk autoincr" json:"menuId"`
	MenuName  string `json:"menuName,omitempty" form:"menuName"`
	ParentId  int64  `json:"parentId,omitempty"`
	OrderNum  int    `json:"orderNum,omitempty"`
	Path      string `json:"path,omitempty"`
	Component string `json:"component,omitempty"`
	Query     string `json:"query,omitempty"`
	IsFrame   string `json:"isFrame,omitempty"`
	IsCache   string `json:"isCache,omitempty"`
	MenuType  string `json:"menuType,omitempty"`
	Visible   string `json:"visible,omitempty" form:"visible"`
	Status    string `json:"status,omitempty" form:"status"`
	Perms     string `json:"perms,omitempty"`
	Icon      string `json:"icon,omitempty"`

	Children []SysMenu `xorm:"-" json:"children,omitempty"`

	// tree select
	Id    int64  `xorm:"-" json:"id,omitempty"`
	Label string `xorm:"-" json:"label,omitempty"`
}
