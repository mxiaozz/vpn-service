package entity

type SysConfig struct {
	BaseEntity `xorm:"extends"`

	ConfigId    int64  `xorm:"pk autoincr" json:"configId" form:"configId"`
	ConfigKey   string `json:"configKey" form:"configKey"`
	ConfigName  string `json:"configName" form:"configName"`
	ConfigType  string `json:"configType" form:"configType"`
	ConfigValue string `json:"configValue" form:"configValue"` // 系统内置（Y是 N否）
}
