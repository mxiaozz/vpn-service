package config

type Token struct {
    Header     string `mapstructure:"header"`     // 令牌自定义标识
    Secret     string `mapstructure:"secret"`     // 令牌秘钥
    Expiration int    `mapstructure:"expiration"` // 令牌有效期（默认30分钟）
}
