package config

type Server struct {
    Host    string `mapstructure:"host" json:"host" yaml:"host"` // IP地址
    Port    string `mapstructure:"port" json:"port" yaml:"port"` // 端口
    DevMode bool   `mapstructure:"dev" json:"dev" yaml:"dev"`    // 是否使用开发模式
}
