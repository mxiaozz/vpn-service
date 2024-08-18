package cst

const (
	CACHE_LOGIN_TOKEN_KEY   = "login_tokens:"  // 登录用户 redis key
	CACHE_CAPTCHA_CODE_KEY  = "captcha_codes:" // 验证码 redis key
	CACHE_SYS_CONFIG_KEY    = "sys_config:"    // 参数管理 cache key
	CACHE_SYS_DICT_KEY      = "sys_dict:"      // 字典管理 cache key
	CACHE_REPEAT_SUBMIT_KEY = "repeat_submit:" // 防重复提交 redis key
	CACHE_RATE_LIMIT_KEY    = "rate_limit:"    // 限流 redis key
	CACHE_PWD_ERR_CNT_KEY   = "pwd_err_cnt:"   // 密码错误次数 redis key
)
