package login

type LoginUser struct {
	// 基础信息
	UserId      int64           `json:"userId"`
	UserName    string          `json:"userName"`
	DeptId      int64           `json:"deptId"`
	DeptName    string          `json:"deptName"`
	Permissions map[string]int8 `json:"permissions,omitempty"`

	// 会话信息
	Token      string `json:"token"`
	ExpireTime int64  `json:"expireTime"`

	// 登录信息
	LoginTime int64  `json:"loginTime"`
	IpAddress string `json:"ipaddr,omitempty"`
	Location  string `json:"location,omitempty"`
	Browser   string `json:"browser,omitempty"`
	AgentOS   string `json:"agentos,omitempty"`
}

func (user LoginUser) IsAdmin() bool {
	return user.UserId == 1
}
