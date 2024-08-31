package config

type FilePath struct {
	Avatar Avatar `mapstructure:"avatar"`
}

type Avatar struct {
	RootStore string `mapstructure:"rootStore"`
	RootUrl   string `mapstructure:"rootUrl"`
}
