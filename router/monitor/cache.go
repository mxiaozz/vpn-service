package monitor

import (
	"vpn-web.funcworks.net/controller/monitor"
	"vpn-web.funcworks.net/router/wraper"
)

func initMonitorCacheRouter(pvt wraper.RouterWraper) {
	cache := wraper.ExtModule("缓存列表")
	{
		pvt.GET("/cache/getNames", monitor.CacheController.GetCacheNames, cache.Ext("monitor:cache:list"))
		pvt.GET("/cache/getKeys/:cacheName", monitor.CacheController.GetCacheKeys, cache.Ext("monitor:cache:list"))
		pvt.GET("/cache/getValue/:cacheName/:cacheKey", monitor.CacheController.GetCacheValue, cache.Ext("monitor:cache:list"))
		pvt.DELETE("/cache/clearCacheName/:cacheName", monitor.CacheController.ClearCacheName, cache.Ext("monitor:cache:list"))
		pvt.DELETE("/cache/clearCacheKey/:cacheKey", monitor.CacheController.ClearCacheKey, cache.Ext("monitor:cache:list"))
		pvt.DELETE("/cache/clearCacheAll", monitor.CacheController.ClearCacheAll, cache.Ext("monitor:cache:list"))
	}
}
