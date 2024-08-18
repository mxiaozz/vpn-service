package response

type SysCache struct {
	CacheName  string `json:"cacheName"`
	CacheKey   string `json:"cacheKey"`
	CacheValue string `json:"cacheValue"`
	Remark     string `json:"remark"`
}
