package config

type Login struct {
    MaxRetryCount int `mapstructure:"maxRetryCount"` // 密码最大错误次数
    LockTime      int `mapstructure:"lockTime"`      // 密码锁定时间（默认10分钟）
}
