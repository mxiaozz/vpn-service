package cst

const (
	USER_STATUS_NORMAL       = "0"          // 用户正常状态
	USER_STATUS_DISABLE      = "1"          // 用户封禁状态
	USER_STATUS_DELETED      = "2"          // 用户删除状态
	USER_USERNAME_MIN_LENGTH = 2            // 用户名最小长度
	USER_USERNAME_MAX_LENGTH = 20           // 用户名最大长度
	USER_PASSWORD_MIN_LENGTH = 5            // 密码最小长度
	USER_PASSWORD_MAX_LENGTH = 20           // 密码最大长度
	ROLE_DISABLE             = "1"          // 角色封禁状态
	DEPT_NORMAL              = "0"          // 部门正常状态
	DEPT_DISABLE             = "1"          // 部门停用状态
	DICT_NORMAL              = "0"          // 字典正常状态
	MENU_YES_FRAME           = "0"          // 是否菜单外链（是）
	MENU_NO_FRAME            = "1"          // 是否菜单外链（否）
	MENU_TYPE_DIR            = "M"          // 菜单类型（目录）
	MENU_TYPE_MENU           = "C"          // 菜单类型（菜单）
	MENU_TYPE_BUTTON         = "F"          // 菜单类型（按钮）
	MENU_LAYOUT              = "Layout"     // Layout组件标识
	MENU_PARENT_VIEW         = "ParentView" // ParentView组件标识
	MENU_INNER_LINK          = "InnerLink"  // InnerLink组件标识
)
