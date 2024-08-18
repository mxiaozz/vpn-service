package model

type UserOnline struct {
	TokenId       string `json:"tokenId"`                  // 会话编号
	DeptName      string `json:"deptName"`                 // 部门名称
	UserName      string `json:"userName" form:"userName"` // 用户名称
	Ipaddr        string `json:"ipaddr" form:"ipaddr"`     // 登录IP地址
	LoginLocation string `json:"loginLocation"`            // 登录地点
	Browser       string `json:"browser"`                  // 浏览器类型
	Os            string `json:"os"`                       // 操作系统
	LoginTime     int64  `json:"loginTime"`                // 登录时间
}
