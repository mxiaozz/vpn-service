package login

import (
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/service/login"
	"vpn-web.funcworks.net/util/rsp"
)

var Captcha = &CaptchaController{}

type CaptchaController struct{}

func (c *CaptchaController) GetCode(ctx *gin.Context) {
	id, b64s, _, err := login.CaptchaService.Generate()
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	data := make(map[string]interface{})
	data["uuid"] = id
	data["img"] = strings.Split(b64s, ";base64,")[1]
	data["captchaEnabled"] = true
	rsp.Context(ctx).Flat().OkWithData(data)
}
