package cst

const (
	SYS_UTF8                 = "UTF-8"          // UTF-8 字符集
	SYS_GBK                  = "GBK"            // GBK 字符集
	SYS_WWW                  = "www."           // www主域
	SYS_HTTP                 = "http://"        // http请求
	SYS_HTTPS                = "https://"       // https请求
	SYS_SUCCESS              = "0"              // 通用成功标识
	SYS_FAIL                 = "1"              // 通用失败标识
	SYS_LOGIN_SUCCESS        = "Success"        // 登录成功
	SYS_LOGOUT               = "Logout"         // 注销
	SYS_REGISTER             = "Register"       // 注册
	SYS_LOGIN_FAIL           = "Error"          // 登录失败
	SYS_ALL_PERMISSION       = "*:*:*"          // 所有权限标识
	SYS_SUPER_ADMIN          = "admin"          // 管理员角色权限标识
	SYS_ROLE_DELIMETER       = ","              // 角色权限分隔符
	SYS_PERMISSION_DELIMETER = ","              // 权限标识分隔符
	SYS_CAPTCHA_EXPIRATION   = 2                // 验证码有效期（分钟）
	SYS_TOKEN                = "token"          // 令牌
	SYS_TOKEN_PREFIX         = "Bearer "        // 令牌前缀
	SYS_LOGIN_USER_KEY       = "login_user_key" // 令牌前缀
	SYS_RESOURCE_PREFIX      = "/profile"       // 资源映射路径 前缀
	SYS_LOOKUP_RMI           = "rmi:"           // RMI 远程方法调用
	SYS_LOOKUP_LDAP          = "ldap:"          // LDAP 远程方法调用
	SYS_LOOKUP_LDAPS         = "ldaps:"         // LDAPS 远程方法调用
	SYS_YES                  = "Y"              // 是否为系统默认（是）
	SYS_UNIQUE               = true             // 唯一
	SYS_NOT_UNIQUE           = false            // 不唯一
)
