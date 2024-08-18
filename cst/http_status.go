package cst

const (
	HTTP_SUCCESS          = 200 // 操作成功
	HTTP_CREATED          = 201 // 对象创建成功
	HTTP_ACCEPTED         = 202 // 请求已经被接受
	HTTP_NO_CONTENT       = 204 // 操作已经执行成功，但是没有返回数据
	HTTP_MOVED_PERM       = 301 // 资源已被移除
	HTTP_SEE_OTHER        = 303 // 重定向
	HTTP_NOT_MODIFIED     = 304 // 资源没有被修改
	HTTP_BAD_REQUEST      = 400 // 参数列表错误（缺少，格式不匹配）
	HTTP_UNAUTHORIZED     = 401 // 未授权
	HTTP_FORBIDDEN        = 403 // 访问受限，授权过期
	HTTP_NOT_FOUND        = 404 // 资源，服务未找到
	HTTP_BAD_METHOD       = 405 // 不允许的请求方式
	HTTP_CONFLICT         = 409 // 资源冲突，或者资源被修改
	HTTP_UNSUPPORTED_TYPE = 415 // 不支持的媒体类型
	HTTP_ERROR            = 500 // 内部错误
	HTTP_NOT_IMPLEMENTED  = 501 // 未实现
	HTTP_WARN             = 601 // 系统警告消息
)
