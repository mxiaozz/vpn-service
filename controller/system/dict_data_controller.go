package system

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var DictDataController = &dictDataController{}

type dictDataController struct {
	controller.BaseController
}

func (c *dictDataController) GetDictDataListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysDictData](ctx)
	if err != nil {
		gb.Logger.Errorln("字典数据列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var dictData entity.SysDictData
	if err = ctx.ShouldBind(&dictData); err != nil {
		gb.Logger.Errorln("字典数据列表获取查询参数失败", err.Error())
		rsp.Fail("字典数据查询参数格式不正确", ctx)
		return
	}

	// 分页查询
	if err = system.DictDataService.GetDictDataListPage(&dictData, page); err != nil {
		gb.Logger.Errorln("字典数据列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *dictDataController) GetDictDataByType(ctx *gin.Context) {
	dictType, exist := ctx.Params.Get("dictType")
	if !exist {
		rsp.Fail("dictType is required", ctx)
		return
	}

	if dictDataList, err := system.DictDataService.GetDictDataByType(dictType); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(dictDataList, ctx)
	}
}

func (c *dictDataController) GetDictData(ctx *gin.Context) {
	dictDataId, _ := strconv.ParseInt(ctx.Param("dictDataId"), 10, 64)
	if dictDataId == 0 {
		gb.Logger.Errorln("获取字典值详情dictCode参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if dictData, err := system.DictDataService.GetDictData(dictDataId); err != nil {
		gb.Logger.Errorln("获取字典值失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(dictData, ctx)
	}
}

func (c *dictDataController) AddDictData(ctx *gin.Context) {
	var dictData entity.SysDictData
	if err := ctx.ShouldBind(&dictData); err != nil {
		gb.Logger.Errorln("增加字典值时，获取字典参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dictData.CreateBy = c.GetLoginUser(ctx).User.UserName
	dictData.CreateTime = time.Now()
	dictData.Status = "0"

	if err := system.DictDataService.AddDictData(&dictData); err != nil {
		gb.Logger.Errorln("增加字典值失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *dictDataController) UpdateDictData(ctx *gin.Context) {
	var dictData entity.SysDictData
	if err := ctx.ShouldBind(&dictData); err != nil {
		gb.Logger.Errorln("修改字典值时，获取字典参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dictData.UpdateBy = c.GetLoginUser(ctx).User.UserName
	dictData.UpdateTime = time.Now()

	if err := system.DictDataService.UpdateDictData(&dictData); err != nil {
		gb.Logger.Errorln("修改字典值失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *dictDataController) DeleteDictData(ctx *gin.Context) {
	dictDataIds := util.NewList(strings.Split(ctx.Param("dictDataIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(dictDataIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.DictDataService.DeleteDictData(dictDataIds); err != nil {
		gb.Logger.Errorln("删除字典值失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
