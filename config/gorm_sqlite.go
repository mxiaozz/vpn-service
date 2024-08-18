package config

type Sqlite struct {
    GeneralDB `yaml:",inline" mapstructure:",squash"`
}
