package config

type StructConfig struct {
    Server  Server          `mapstructure:"server" json:"server" yaml:"server"`
    Token   Token           `mapstructure:"token" json:"token" yaml:"token"`
    Redis   Redis           `mapstructure:"redis" json:"redis" yaml:"redis"`
    Captcha Captcha         `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
    Sqlite  Sqlite          `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
    DBList  []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
    Login   Login           `mapstructure:"login" json:"login" yaml:"login"`
}
