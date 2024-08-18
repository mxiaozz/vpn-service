package entity

// 岗位表 sys_post
type SysPost struct {
	BaseEntity `xorm:"extends"`

	PostId   int64  `xorm:"pk autoincr" json:"postId"` // 岗位ID
	PostCode string `json:"postCode" form:"postCode"`  // 岗位编码
	PostName string `json:"postName" form:"postName"`  // 岗位名称
	PostSort int    `json:"postSort"`                  // 岗位排序
	Status   string `json:"status" form:"status"`      // 状态（0正常 1停用）
}
