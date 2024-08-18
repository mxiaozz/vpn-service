package monitor

import (
	"context"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/response"
	"vpn-web.funcworks.net/util/rsp"
)

var CacheController = &cacheController{}

var cacheNames = []response.SysCache{
	{CacheName: cst.CACHE_LOGIN_TOKEN_KEY, Remark: "用户信息"},
	{CacheName: cst.CACHE_SYS_CONFIG_KEY, Remark: "配置信息"},
	{CacheName: cst.CACHE_SYS_DICT_KEY, Remark: "数据字典"},
	{CacheName: cst.CACHE_CAPTCHA_CODE_KEY, Remark: "验证码"},
	{CacheName: cst.CACHE_PWD_ERR_CNT_KEY, Remark: "密码错误次数"},
}

type cacheController struct {
	controller.BaseController
}

func (c *cacheController) GetCacheNames(ctx *gin.Context) {
	rsp.OkWithData(cacheNames, ctx)
}

func (c *cacheController) GetCacheKeys(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	if cacheName == "" {
		rsp.Fail("参数错误", ctx)
		return
	}

	if keys, err := gb.RedisClient.Keys(context.Background(), cacheName+"*").Result(); err != nil {
		gb.Logger.Error(err)
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		rsp.OkWithData(keys, ctx)
	}
}

func (c *cacheController) GetCacheValue(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	cacheKey := ctx.Param("cacheKey")
	if cacheKey == "" {
		rsp.Fail("参数错误", ctx)
		return
	}

	if value, err := gb.RedisProxy.Get(cacheKey); err != nil {
		gb.Logger.Error(err)
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		rsp.OkWithData(response.SysCache{
			CacheName:  cacheName,
			CacheKey:   cacheKey,
			CacheValue: value},
			ctx)
	}
}

func (c *cacheController) ClearCacheName(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	if cacheName == "" {
		rsp.Fail("参数错误", ctx)
		return
	}

	if keys, err := gb.RedisClient.Keys(context.Background(), cacheName+"*").Result(); err != nil {
		gb.Logger.Error(err)
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		for _, k := range keys {
			if err := gb.RedisProxy.Delete(k); err != nil {
				gb.Logger.Error(err)
				rsp.Fail(err.Error(), ctx)
				return
			}
		}
	}
	rsp.Ok(ctx)
}

func (c *cacheController) ClearCacheKey(ctx *gin.Context) {
	cacheKey := ctx.Param("cacheKey")
	if cacheKey == "" {
		rsp.Fail("参数错误", ctx)
		return
	}

	if err := gb.RedisProxy.Delete(cacheKey); err != nil {
		gb.Logger.Error(err)
		rsp.Fail(err.Error(), ctx)
		return
	}
	rsp.Ok(ctx)
}

func (c *cacheController) ClearCacheAll(ctx *gin.Context) {
	if keys, err := gb.RedisClient.Keys(context.Background(), "*").Result(); err != nil {
		gb.Logger.Error(err)
		rsp.Fail(err.Error(), ctx)
		return
	} else {
		for _, k := range keys {
			if err := gb.RedisProxy.Delete(k); err != nil {
				gb.Logger.Error(err)
				rsp.Fail(err.Error(), ctx)
				return
			}
		}
	}
	rsp.Ok(ctx)
}
