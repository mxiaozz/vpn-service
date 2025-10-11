package controller

import (
	"io"
	"slices"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/gb"
	service "vpn-web.funcworks.net/service/openvpn"
	"vpn-web.funcworks.net/util/rsp"
)

var Openvpn = &openvpnController{}

type openvpnController struct {
	BaseController
}

func (c *openvpnController) GetStatus(ctx *gin.Context) {
	status := service.OpenvpnService.GetServerStatus()
	rsp.OkWithData(status, ctx)
}

func (c *openvpnController) GetRealtimeStatus(ctx *gin.Context) {
	status := service.OpenvpnService.GetRealtimeStatus()
	rsp.OkWithData(status, ctx)
}

func (c *openvpnController) GetServerConfig(ctx *gin.Context) {
	if cfg, err := service.OpenvpnService.GetServerConfig(); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(cfg, ctx)
	}
}

func (c *openvpnController) GenerateConfig(ctx *gin.Context) {
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}
	if cfg, err := service.OpenvpnService.GenerateConfig(string(data)); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(cfg, ctx)
	}
}

func (c *openvpnController) SaveConfig(ctx *gin.Context) {
	var data map[string]string
	if err := ctx.ShouldBindJSON(&data); err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	cfgContent := data["cfgContent"]
	if cfgContent == "" {
		rsp.Fail("cfgContent is required", ctx)
		return
	}

	if cfg, err := service.OpenvpnService.SaveConfig(cfgContent); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(cfg, ctx)
	}
}

func (c *openvpnController) OptServer(ctx *gin.Context) {
	opt, exist := ctx.GetQuery("opt")
	if !exist {
		rsp.Fail("opt is required", ctx)
		return
	}

	list := []string{"start", "stop", "restart"}
	if !slices.Contains(list, opt) {
		rsp.Fail("opt is error", ctx)
		return
	}

	service.OpenvpnService.ChangeServer(string(opt))
	rsp.Ok(ctx)
}

func (c *openvpnController) GetPKIStatus(ctx *gin.Context) {
	if isInited, err := service.OpenvpnService.GetPKIStatus(); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(isInited, ctx)
	}
}

func (c *openvpnController) InikPKI(ctx *gin.Context) {
	if err := service.OpenvpnService.InikPKI(); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *openvpnController) ResetPKI(ctx *gin.Context) {
	if err := service.OpenvpnService.ResetPKI(); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *openvpnController) GetUserCert(ctx *gin.Context) {
	userName, exist := ctx.GetQuery("userName")
	if !exist {
		rsp.Fail("userName is required", ctx)
		return
	}

	if userCert, err := service.OpenvpnService.GetUserCert(userName, false); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(userCert, ctx)
	}
}

func (c *openvpnController) GenerateUserCert(ctx *gin.Context) {
	userName, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}
	if err = service.OpenvpnService.GenerateUserCert(string(userName)); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *openvpnController) RevokeUserCert(ctx *gin.Context) {
	userName, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}
	if err = service.OpenvpnService.RevokeUserCert(string(userName)); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *openvpnController) DownUserCert(ctx *gin.Context) {
	userName, _ := ctx.GetPostForm("userName")
	gb.Logger.Infof("begin download user: %s cert", userName)
	if userName == "" {
		rsp.Fail("userName is required", ctx)
		return
	}

	if userCert, err := service.OpenvpnService.GetUserCert(userName, true); err != nil {
		rsp.Fail(err.Error(), ctx)
	} else {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteString(userCert.Cert)
		ctx.Writer.Flush()
	}
}
