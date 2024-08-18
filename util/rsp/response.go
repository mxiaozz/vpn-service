package rsp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
)

type response struct {
	ctx  *gin.Context
	flat bool
}

func Context(ctx *gin.Context) *response {
	return &response{ctx: ctx}
}

func (r *response) Context(ctx *gin.Context) *response {
	r.ctx = ctx
	return r
}

func Flat() *response {
	return &response{flat: true}
}

func (r *response) Flat() *response {
	r.flat = true
	return r
}

func Ok(ctx *gin.Context) {
	write(cst.HTTP_SUCCESS, nil, "", false, ctx)
}
func OkWithMessage(msg string, ctx *gin.Context) {
	write(cst.HTTP_SUCCESS, nil, msg, false, ctx)
}

func OkWithData(data any, ctx *gin.Context) {
	write(cst.HTTP_SUCCESS, data, "", false, ctx)
}

func OkWithResponse(data *model.Response[any], ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func (r *response) OkWithData(data any) {
	write(cst.HTTP_SUCCESS, data, "", r.flat, r.ctx)
}

func Fail(msg string, ctx *gin.Context) {
	write(cst.HTTP_ERROR, nil, msg, false, ctx)
}

func FailWithCode(code int, msg string, ctx *gin.Context) {
	write(code, nil, msg, false, ctx)
}

func write(code int, obj any, msg string, flat bool, ctx *gin.Context) {
	if ctx == nil {
		gb.Logger.Errorf("gin.Context is nil")
		return
	}

	data := make(map[string]any)

	if obj != nil {
		if flat {
			if m, ok := obj.(map[string]any); ok {
				for k, v := range m {
					data[k] = v
				}
				goto to
			}
		}
		data["data"] = obj
	}
to:
	if msg != "" {
		data["msg"] = msg
	}
	data["code"] = code

	ctx.JSON(http.StatusOK, data)
}
