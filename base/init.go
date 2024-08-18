package base

func init() {
	// 初始化日志
	initLogger()

	// 设置默认配置，在配置加载之前
	initDefaultConfig()

	// 加载配置
	initYamlConfig()

	// 初始化数据库
	initSqlite()

	// 初始化 redis client
	initRedis()

	// 初始化定时任务
	initSched()
}
