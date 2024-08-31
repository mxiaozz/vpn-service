package config

type StructConfig struct {
	Server   Server          `mapstructure:"server"`
	Token    Token           `mapstructure:"token"`
	Redis    Redis           `mapstructure:"redis"`
	Captcha  Captcha         `mapstructure:"captcha"`
	Sqlite   Sqlite          `mapstructure:"sqlite"`
	DBList   []SpecializedDB `mapstructure:"db-list"`
	Login    Login           `mapstructure:"login"`
	FilePath FilePath        `mapstructure:"file"`
}
