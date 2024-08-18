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

var DictController = &dictController{}

type dictController struct {
	controller.BaseController
}

func (c *dictController) GetAllDicts(ctx *gin.Context) {
	if dicts, err := system.DictService.GetAllDicts(); err != nil {
		gb.Logger.Errorln("字典列表获取字典列表失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(dicts, ctx)
	}
}

func (c *dictController) GetDictListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysDictType](ctx)
	if err != nil {
		gb.Logger.Errorln("字典列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var dict entity.SysDictType
	if err = ctx.ShouldBind(&dict); err != nil {
		gb.Logger.Errorln("字典列表获取查询参数失败", err.Error())
		rsp.Fail("字典查询参数格式不正确", ctx)
		return
	}
	if dict.Params == nil {
		dict.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			dict.Params[k] = v
		}
	}

	// 分页查询
	if err = system.DictService.GetDictListPage(&dict, page); err != nil {
		gb.Logger.Errorln("字典列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *dictController) GetDict(ctx *gin.Context) {
	dictId, _ := strconv.ParseInt(ctx.Param("dictId"), 10, 64)
	if dictId == 0 {
		gb.Logger.Errorln("字典列表获取字典详情dictId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if dict, err := system.DictService.GetDict(dictId); err != nil {
		gb.Logger.Errorln("字典列表获取岗位详情失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(dict, ctx)
	}
}

func (c *dictController) AddDict(ctx *gin.Context) {
	var dict entity.SysDictType
	if err := ctx.ShouldBind(&dict); err != nil {
		gb.Logger.Errorln("增加字典时，获取字典参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dict.CreateBy = c.GetLoginUser(ctx).User.UserName
	dict.CreateTime = time.Now()
	dict.Status = "0"

	if err := system.DictService.AddDict(&dict); err != nil {
		gb.Logger.Errorln("增加字典失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *dictController) UpdateDict(ctx *gin.Context) {
	var dict entity.SysDictType
	if err := ctx.ShouldBind(&dict); err != nil {
		gb.Logger.Errorln("修改字典时，获取字典参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	dict.UpdateBy = c.GetLoginUser(ctx).User.UserName
	dict.UpdateTime = time.Now()

	if err := system.DictService.UpdateDict(&dict); err != nil {
		gb.Logger.Errorln("修改字典失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *dictController) DeleteDict(ctx *gin.Context) {
	dictIds := util.NewList(strings.Split(ctx.Param("dictIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(dictIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.DictService.DeleteDict(dictIds); err != nil {
		gb.Logger.Errorln("删除字典失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *dictController) ReloadDictCache(ctx *gin.Context) {
	if err := system.DictService.ReloadConfigCache(); err != nil {
		gb.Logger.Errorln("刷新缓存失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
