package entity

type SysDictData struct {
	BaseEntity `xorm:"extends"`

	DictCode  int64  `xorm:"pk autoincr" json:"dictCode"` // 字典编码
	DictSort  int64  `json:"dictSort"`                    // 字典排序
	DictLabel string `json:"dictLabel" form:"dictLabel"`  // 字典标签
	DictValue string `json:"dictValue"`                   // 字典键值
	DictType  string `json:"dictType" form:"dictType"`    // 字典类型
	CssClass  string `json:"cssClass"`                    // 样式属性（其他样式扩展）
	ListClass string `json:"listClass"`                   // 表格字典样式
	IsDefault string `json:"isDefault"`                   // 是否默认（Y是 N否）
	Status    string `json:"status" form:"status"`        // 状态（0正常 1停用）
}
