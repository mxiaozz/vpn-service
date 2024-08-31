package system

import (
	"encoding/json"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
)

var DictDataService = &dictDataService{}

type dictDataService struct {
}

func (ds *dictDataService) GetDictDataListPage(dictData entity.SysDictData, page *model.Page[entity.SysDictData]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_dict_data").
			Where(builder.If(dictData.DictType != "", builder.Eq{"dict_type": dictData.DictType}).
				And(builder.If(dictData.DictLabel != "", builder.Like{"dict_label", dictData.DictLabel})).
				And(builder.If(dictData.Status != "", builder.Eq{"status": dictData.Status})))
		return builder.Expr("dict_sort")
	})
}

func (ds *dictDataService) GetDictDataByType(dictType string) ([]entity.SysDictData, error) {
	var dictData = []entity.SysDictData{}
	// 读取缓存
	if dataStr, err := gb.RedisProxy.Get(ds.getCacheKey(dictType)); err == nil {
		if err := json.Unmarshal([]byte(dataStr), &dictData); err == nil {
			return dictData, nil
		}
	}
	// 数据库读取并缓存
	if err := gb.DB.Where("status = '0' and dict_type = ?", dictType).Asc("dict_sort").Find(&dictData); err != nil {
		return nil, err
	} else {
		if dictBytes, err := json.Marshal(dictData); err != nil {
			gb.Logger.Warn(err.Error())
		} else if err := gb.RedisProxy.Set(ds.getCacheKey(dictType), string(dictBytes)); err != nil {
			gb.Logger.Warn(err.Error())
		}
		return dictData, err
	}
}

func (ds *dictDataService) GetDictData(dictDataId int64) (entity.SysDictData, error) {
	var dictData entity.SysDictData
	if exist, err := gb.DB.Where("dict_code = ?", dictDataId).Get(&dictData); err != nil {
		return dictData, err
	} else if !exist {
		return dictData, errors.Wrap(gb.ErrNotFound, "字典数据不存在")
	}
	return dictData, nil
}

func (ds *dictDataService) AddDictData(dictData entity.SysDictData) error {
	if exist, err := ds.checkDictValueUnique(dictData.DictType, dictData.DictValue, dictData.DictCode); err != nil {
		return err
	} else if exist {
		return errors.New("字典值已存在")
	}

	_, err := gb.DB.Insert(dictData)
	return err
}

func (ds *dictDataService) UpdateDictData(dictData entity.SysDictData) error {
	if exist, err := ds.checkDictValueUnique(dictData.DictType, dictData.DictValue, dictData.DictCode); err != nil {
		return err
	} else if exist {
		return errors.New("字典值已存在")
	}

	if _, err := gb.DB.Where("dict_code = ?", dictData.DictCode).Update(dictData); err != nil {
		return err
	} else {
		return gb.RedisProxy.Delete(ds.getCacheKey(dictData.DictType))
	}
}

func (ds *dictDataService) checkDictValueUnique(dictType, dictValue string, dictDataId int64) (bool, error) {
	return gb.DB.Table("sys_dict_data").Where("dict_type = ? and dict_value = ? and dict_code <> ?", dictType, dictValue, dictDataId).Exist()
}

func (ds *dictDataService) DeleteDictData(dictDataIds []int64) error {
	var dtList = []string{}
	if err := gb.DB.Table("sys_dict_data").Select("dict_type").Distinct().In("dict_code", dictDataIds).Find(&dtList); err != nil {
		return err
	}
	for _, t := range dtList {
		if err := gb.RedisProxy.Delete(ds.getCacheKey(t)); err != nil {
			return err
		}
	}
	_, err := gb.DB.Table("sys_dict_data").In("dict_code", dictDataIds).Delete()
	return err
}

func (ds *dictDataService) getCacheKey(dictType string) string {
	return cst.CACHE_SYS_DICT_KEY + dictType
}
