package system

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var DictService = &dictService{}

type dictService struct {
}

func (ds *dictService) GetAllDicts() ([]entity.SysDictType, error) {
	var dicts = []entity.SysDictType{}
	err := gb.DB.OrderBy("dict_id").Find(&dicts)
	return dicts, err
}

func (ds *dictService) GetDictListPage(dict entity.SysDictType, page *model.Page[entity.SysDictType]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_dict_type").
			Where(builder.If(dict.DictName != "", builder.Like{"dict_name", dict.DictName}).
				And(builder.If(dict.DictType != "", builder.Like{"dict_type", dict.DictType})).
				And(builder.If(dict.Status != "", builder.Eq{"status": dict.Status})).
				And(builder.If(func() bool { return dict.Params["beginTime"] != nil }(),
					builder.Gte{"create_time": dict.Params["beginTime"]})).
				And(builder.If(func() bool { return dict.Params["endTime"] != nil }(),
					builder.Lte{"create_time": dict.Params["endTime"]})))
		return builder.Expr("dict_id")
	})
}

func (ds *dictService) GetDict(dictId int64) (entity.SysDictType, error) {
	var dict entity.SysDictType
	if exist, err := gb.DB.Where("dict_id = ?", dictId).Get(&dict); err != nil {
		return dict, err
	} else if !exist {
		return dict, errors.Wrap(gb.ErrNotFound, "字典不存在")
	}
	return dict, nil
}

func (ds *dictService) AddDict(dict entity.SysDictType) error {
	if exist, err := ds.checkDictTypeUnique(dict.DictType, dict.DictId); err != nil {
		return err
	} else if exist {
		return errors.New("字典类型已存在")
	}

	_, err := gb.DB.Insert(dict)
	return err
}

func (ds *dictService) UpdateDict(dict entity.SysDictType) error {
	if exist, err := ds.checkDictTypeUnique(dict.DictType, dict.DictId); err != nil {
		return err
	} else if exist {
		return errors.New("字典类型已存在")
	}

	_, err := gb.DB.Where("dict_id = ?", dict.DictId).Update(dict)
	return err
}

func (ds *dictService) checkDictTypeUnique(dictType string, dictId int64) (bool, error) {
	return gb.DB.Table("sys_dict_type").Where("dict_type = ? and dict_id <> ?", dictType, dictId).Exist()
}

func (ds *dictService) DeleteDict(dictIds []int64) error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		for _, dictId := range dictIds {
			dict, err := ds.GetDict(dictId)
			if err == gb.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}

			if exist, err := ds.checkDictTypeData(dict.DictType); err != nil {
				return err
			} else if exist {
				return errors.New(fmt.Sprintf("字典[%s]存在具体数据，不能删除", dict.DictName))
			}

			// 删除缓存
			if err := gb.RedisProxy.Delete(ds.getCacheKey(dict.DictType)); err != nil {
				return err
			}

			// 删除字典数据
			if _, err := gb.DB.Table("sys_dict_type").Where("dict_id = ?", dictId).Delete(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ds *dictService) checkDictTypeData(dictType string) (bool, error) {
	return gb.DB.Table("sys_dict_data").Where("dict_type = ?", dictType).Exist()
}

func (ds *dictService) ReloadConfigCache() error {
	cacheKey := ds.getCacheKey("*")
	if err := gb.RedisProxy.Delete(cacheKey); err != nil {
		return err
	}

	var dicts = []entity.SysDictType{}
	if err := gb.DB.Find(&dicts); err != nil {
		return err
	}

	for _, dict := range dicts {
		if dictList, err := DictDataService.GetDictDataByType(dict.DictType); err != nil {
			return err
		} else if dictBytes, err := json.Marshal(dictList); err != nil {
			return err
		} else {
			if err := gb.RedisProxy.Set(ds.getCacheKey(dict.DictType), string(dictBytes)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ds *dictService) getCacheKey(dictType string) string {
	return cst.CACHE_SYS_DICT_KEY + dictType
}
