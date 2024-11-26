package security

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/model/login"
	"vpn-web.funcworks.net/model/request"
	"vpn-web.funcworks.net/service/monitor"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
)

type logContext struct {
	RequestBody  []byte
	ResponseBody []byte

	Module   string
	Path     string
	Method   string
	ClientId string
	Browser  string
	Os       string

	BeginTime model.DateTime
	CostTime  time.Duration

	UserId   int64
	UserName string
	DeptName string

	ResponseBuffer      *bytes.Buffer
	ResponseContentType string
}

func saveLog(ctx *logContext) {
	parts := strings.Split(ctx.Path, "/")
	lastPart := strings.ToLower(parts[len(parts)-1]) // 转换为小写方便一致性比较
	if lastPart == "login" {
		setLoginUserName(ctx)
		saveLoginLog(ctx)
	} else if lastPart == "logout" {
		saveLoginLog(ctx)
	} else {
		if ctx.Method == http.MethodPost && lastPart == "user" ||
			ctx.Method == http.MethodPut && lastPart == "resetpwd" {
			// 用户管理：新增用户 | 重置密码
			var data map[string]any
			if err := json.Unmarshal(ctx.RequestBody, &data); err != nil {
				gb.Logger.Errorln("保存日志", "解析请求json失败", err.Error())
				ctx.RequestBody = nil
			}
			data["password"] = "***"
			ctx.RequestBody, _ = json.Marshal(data)
		} else if ctx.Method == http.MethodPut && lastPart == "updatepwd" {
			// 个人中心：修改密码
			var data map[string]any
			if err := json.Unmarshal(ctx.RequestBody, &data); err != nil {
				gb.Logger.Errorln("保存日志", "解析请求json失败", err.Error())
				ctx.RequestBody = nil
			}
			data["oldPassword"] = "***"
			data["newPassword"] = "***"
			ctx.RequestBody, _ = json.Marshal(data)
		}
		saveOperLog(ctx)
	}
}

func setLoginUserName(ctx *logContext) {
	var loginRequest request.LoginRequest
	if err := json.Unmarshal(ctx.RequestBody, &loginRequest); err != nil {
		gb.Logger.Errorln("保存登录日志", "解析请求json失败", err.Error())
		return
	}
	ctx.UserName = loginRequest.Username
}

func saveLoginLog(ctx *logContext) {
	// 解析返回值
	code, msg := parseResponseBody(ctx, "登录成功")

	// 保存日志
	log := entity.SysLoginLog{
		UserName:      ctx.UserName,
		Ipaddr:        ctx.ClientId,
		LoginLocation: "",
		Browser:       ctx.Browser,
		Os:            ctx.Os,
		Msg:           msg,
		LoginTime:     ctx.BeginTime,
		Status:        util.If(code == cst.HTTP_SUCCESS, cst.SYS_SUCCESS, cst.SYS_FAIL),
	}

	if err := monitor.LoginLogService.AddLoginLog(log); err != nil {
		gb.Logger.Errorln("保存登录日志失败", err.Error())
	}
}

func parseResponseBody(ctx *logContext, defaultMsg string) (code int, msg string) {
	var data map[string]any
	if strings.HasPrefix(ctx.ResponseContentType, "application/json") {
		if err := json.Unmarshal(ctx.ResponseBody, &data); err != nil {
			gb.Logger.Errorln("保存日志", "解析返json数失败", err.Error())
			return
		}
	}
	code, msg = cst.HTTP_SUCCESS, defaultMsg
	if v, ok := data["code"]; ok {
		code = int(v.(float64))
	}
	if v, ok := data["msg"]; ok {
		msg = v.(string)
	}
	return
}

func saveOperLog(ctx *logContext) {
	// 解析返回值
	code, msg := parseResponseBody(ctx, "")

	// 保存日志
	operLog := entity.SysOperLog{
		Title:         ctx.Module,
		BusinessType:  0,
		Method:        "",
		RequestMethod: ctx.Method,
		OperatorType:  1,
		OperName:      ctx.UserName,
		DeptName:      ctx.DeptName,
		OperUrl:       ctx.Path,
		OperIp:        ctx.ClientId,
		OperLocation:  "",
		ErrorMsg:      msg,
		CostTime:      ctx.CostTime.Milliseconds(),
		OperTime:      ctx.BeginTime,
		Status:        util.If(code == cst.HTTP_SUCCESS, 0, 1),
	}
	if len(ctx.RequestBody) > 1024 {
		operLog.OperParam = string(ctx.RequestBody[:1024])
	} else {
		operLog.OperParam = string(ctx.RequestBody)
	}
	if len(ctx.ResponseBody) > 1024 {
		operLog.JsonResult = string(ctx.ResponseBody[:1024])
	} else {
		operLog.JsonResult = string(ctx.ResponseBody)
	}

	switch ctx.Method {
	case "POST":
		operLog.BusinessType = 1
	case "PUT":
		operLog.BusinessType = 2
	case "DELETE":
		operLog.BusinessType = 3
	}

	if err := system.OperLogService.AddOperLog(operLog); err != nil {
		gb.Logger.Errorln("保存操作日志失败", err.Error())
	}
}

func doBefore(ext ExtInfo, ctx *gin.Context) (*logContext, error) {
	// GET 操作不记录日志
	if strings.EqualFold(ctx.Request.Method, "GET") {
		return nil, nil
	}

	// request body
	body, err := getRequestBody(ctx)
	if err != nil {
		return nil, err
	}
	logCtx := &logContext{
		Module:      ext.Module,
		RequestBody: body,
		Path:        ctx.Request.URL.Path,
		Method:      ctx.Request.Method,
		ClientId:    ctx.ClientIP(),
	}

	// parse user agent
	ua := useragent.New(ctx.GetHeader("User-Agent"))
	logCtx.Browser, _ = ua.Browser()
	logCtx.Os = ua.OS()

	// login user info
	if user, _ := ctx.Get(cst.SYS_LOGIN_USER_KEY); user != nil {
		loginUser := user.(login.LoginUser)
		logCtx.UserId = loginUser.UserId
		logCtx.UserName = loginUser.UserName
		logCtx.DeptName = loginUser.DeptName
	}

	// 接收响应 body
	writer := responseWriter{
		ResponseWriter: ctx.Writer,
		body:           &bytes.Buffer{},
	}
	ctx.Writer = writer
	logCtx.ResponseBuffer = writer.body

	// 业务耗时记时开始
	logCtx.BeginTime = model.DateTimeNow()

	return logCtx, nil
}

func doAfter(logCtx *logContext, ctx *gin.Context) {
	if logCtx == nil {
		return
	}

	logCtx.CostTime = time.Since(logCtx.BeginTime.Time)
	logCtx.ResponseBody = logCtx.ResponseBuffer.Bytes()
	logCtx.ResponseBuffer = nil
	logCtx.ResponseContentType = ctx.Writer.Header().Get("Content-Type")

	go saveLog(logCtx)
}

func getRequestBody(ctx *gin.Context) ([]byte, error) {
	if ctx.Request.ContentLength <= 0 {
		if len(ctx.Request.URL.RawQuery) > 0 {
			return []byte(ctx.Request.URL.RawQuery), nil
		}
		return nil, nil
	}

	if strings.Contains(ctx.GetHeader("Content-Type"), "multipart/form-data") {
		return []byte("[文件]"), nil
	}

	if body, err := io.ReadAll(ctx.Request.Body); err != nil {
		return nil, err
	} else {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return body, nil
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
