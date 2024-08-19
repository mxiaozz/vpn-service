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

var ConfigController = &configController{}

type configController struct {
	controller.BaseController
}

func (c *configController) GetConfigListPage(ctx *gin.Context) {
	// 获取分页参数
	page, err := model.NewPage[entity.SysConfig](ctx)
	if err != nil {
		gb.Logger.Errorln("参数列表获取分页参数失败", err.Error())
		rsp.Fail("获取分页参数失败", ctx)
		return
	}

	// 获取查询参数
	var config entity.SysConfig
	if err = ctx.ShouldBind(&config); err != nil {
		gb.Logger.Errorln("参数列表获取查询参数失败", err.Error())
		rsp.Fail("参数查询参数格式不正确", ctx)
		return
	}
	if config.Params == nil {
		config.Params = make(map[string]any)
		params := ctx.QueryMap("params")
		for k, v := range params {
			config.Params[k] = v
		}
	}

	// 分页查询
	if err = system.ConfigService.GetConfigListPage(&config, page); err != nil {
		gb.Logger.Errorln("角色列表查询失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(page.ToMap())
	}
}

func (c *configController) GetConfigByKey(ctx *gin.Context) {
	configKey, exist := ctx.Params.Get("configKey")
	if !exist {
		rsp.Fail("configKey is required", ctx)
		return
	}

	if configValue, err := system.ConfigService.GetConfigByKey(configKey); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(configValue, ctx)
	}
}

func (c *configController) GetConfig(ctx *gin.Context) {
	configId, _ := strconv.ParseInt(ctx.Param("configId"), 10, 64)
	if configId == 0 {
		gb.Logger.Errorln("获取参数详情configId参数错误")
		rsp.Fail("参数错误", ctx)
		return
	}

	if config, err := system.ConfigService.GetConfig(configId); err != nil {
		gb.Logger.Errorln("获取参数失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(config, ctx)
	}
}

func (c *configController) AddConfig(ctx *gin.Context) {
	var config entity.SysConfig
	if err := ctx.ShouldBind(&config); err != nil {
		gb.Logger.Errorln("获取参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	config.CreateBy = c.GetLoginUser(ctx).User.UserName
	config.CreateTime = model.DateTime(time.Now())

	if err := system.ConfigService.AddConfig(&config); err != nil {
		gb.Logger.Errorln("增加参数失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *configController) UpdateConfig(ctx *gin.Context) {
	var config entity.SysConfig
	if err := ctx.ShouldBind(&config); err != nil {
		gb.Logger.Errorln("获取参数对象失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 增补信息
	config.UpdateBy = c.GetLoginUser(ctx).User.UserName
	config.UpdateTime = model.DateTime(time.Now())

	if err := system.ConfigService.UpdateConfig(&config); err != nil {
		gb.Logger.Errorln("修参数失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *configController) DeleteConfig(ctx *gin.Context) {
	configIds := util.NewList(strings.Split(ctx.Param("configIds"), ",")).
		Filter(func(id string) bool { return id != "" }).
		Distinct(func(id string) any { return id }).
		MapToInt64(func(id string) int64 {
			uid, _ := strconv.ParseInt(id, 10, 64)
			return uid
		}).
		Filter(func(id int64) bool { return id > 0 })
	if len(configIds) == 0 {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := system.ConfigService.DeleteConfig(configIds); err != nil {
		gb.Logger.Errorln("删除参数失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *configController) ReloadConfigCache(ctx *gin.Context) {
	if err := system.ConfigService.ReloadConfigCache(); err != nil {
		gb.Logger.Errorln("刷新缓存失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}
