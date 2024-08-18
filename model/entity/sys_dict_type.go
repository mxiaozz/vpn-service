package entity

type SysDictType struct {
	BaseEntity `xorm:"extends"`

	DictId   int64  `xorm:"pk autoincr" json:"dictId"` // 字典主键
	DictName string `json:"dictName" form:"dictName"`  // 字典名称
	DictType string `json:"dictType" form:"dictType"`  // 字典类型
	Status   string `json:"status" form:"status"`      // 状态（0正常 1停用）
}
