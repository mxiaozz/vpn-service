package system

import (
	"github.com/pkg/errors"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
)

var ConfigService = &configService{}

type configService struct {
}

// 分页查询
func (cs *configService) GetConfigListPage(config *entity.SysConfig, page *model.Page[entity.SysConfig]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_config").
			Where(builder.If(config.ConfigName != "", builder.Like{"config_name", config.ConfigName}).
				And(builder.If(config.ConfigType != "", builder.Eq{"config_type": config.ConfigType})).
				And(builder.If(config.ConfigKey != "", builder.Like{"config_key", config.ConfigKey})).
				And(builder.If(func() bool { return config.Params["beginTime"] != nil }(),
					builder.Gte{"create_time": config.Params["beginTime"]})).
				And(builder.If(func() bool { return config.Params["endTime"] != nil }(),
					builder.Lte{"create_time": config.Params["endTime"]})))
		return builder.Expr("config_id")
	})
}

// 根据参数键名查询参数值，先查缓存，无缓存时查数据库
func (cs *configService) GetConfigByKey(configKey string) (string, error) {
	cacheKey := cs.getCacheKey(configKey)
	if value, err := gb.RedisProxy.Get(cacheKey); err == nil {
		return value, nil
	}

	var config entity.SysConfig
	if _, err := gb.DB.Where("config_key = ?", configKey).Get(&config); err == nil {
		if err := gb.RedisProxy.Set(cacheKey, config.ConfigValue); err != nil {
			gb.Logger.Errorln("设置参数缓存失败", err)
		}
		return config.ConfigValue, nil
	} else {
		return "", err
	}
}

func (cs *configService) GetConfig(configId int64) (*entity.SysConfig, error) {
	var config entity.SysConfig
	if exist, err := gb.DB.Where("config_id = ?", configId).Get(&config); err != nil || !exist {
		return nil, err
	}
	return &config, nil
}

func (cs *configService) AddConfig(config *entity.SysConfig) error {
	if exist, err := cs.checkConfigKeyUnique(config.ConfigKey, config.ConfigId); err != nil {
		return err
	} else if exist {
		return errors.New("参数已存在")
	}

	_, err := gb.DB.Insert(config)
	return err
}

func (cs *configService) UpdateConfig(config *entity.SysConfig) error {
	if exist, err := cs.checkConfigKeyUnique(config.ConfigKey, config.ConfigId); err != nil {
		return err
	} else if exist {
		return errors.New("参数已存在")
	}

	if _, err := gb.DB.Where("config_id = ?", config.ConfigId).Update(config); err != nil {
		return err
	}
	return gb.RedisProxy.Delete(cs.getCacheKey(config.ConfigKey))
}

func (cs *configService) checkConfigKeyUnique(configKey string, configId int64) (bool, error) {
	return gb.DB.Table("sys_config").Where("config_key = ? and config_id <> ?", configKey, configId).Exist()
}

func (cs *configService) DeleteConfig(configIds []int64) error {
	// 删除缓存
	var configs = []entity.SysConfig{}
	if err := gb.DB.In("config_id", configIds).Find(&configs); err != nil {
		return err
	}
	for _, config := range configs {
		if err := gb.RedisProxy.Delete(cs.getCacheKey(config.ConfigKey)); err != nil {
			return err
		}
	}

	// 删除数据
	_, err := gb.DB.Table("sys_config").In("config_id", configIds).Delete()
	return err
}

func (cs *configService) ReloadConfigCache() error {
	cacheKey := cs.getCacheKey("*")
	if err := gb.RedisProxy.Delete(cacheKey); err != nil {
		return err
	}

	var configs = []entity.SysConfig{}
	if err := gb.DB.Find(&configs); err != nil {
		return err
	}

	for _, config := range configs {
		if err := gb.RedisProxy.Set(cs.getCacheKey(config.ConfigKey), config.ConfigValue); err != nil {
			return err
		}
	}
	return nil
}

func (cs *configService) getCacheKey(configKey string) string {
	return cst.CACHE_SYS_CONFIG_KEY + configKey
}
