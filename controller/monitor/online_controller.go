package monitor

import (
	"cmp"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/model/login"
	"vpn-web.funcworks.net/model/response"
	"vpn-web.funcworks.net/service/openvpn"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var OnlineController = &onlineController{}

type onlineController struct {
	controller.BaseController
}

func (c *onlineController) GetOnlineUsers(ctx *gin.Context) {
	// 获取查询条件
	var userOnline response.UserOnline
	if err := ctx.ShouldBind(&userOnline); err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 从缓存中获取在线用户
	keys, err := gb.RedisClient.Keys(context.Background(), cst.CACHE_LOGIN_TOKEN_KEY+"*").Result()
	if err != nil {
		gb.Logger.Errorln("缓存获取在线用户失败", err.Error())
		rsp.Fail("获取在线用户失败", ctx)
		return
	}

	var users = make([]response.UserOnline, 0)
	for _, key := range keys {
		userStr, err := gb.RedisProxy.Get(key)
		if err != nil {
			gb.Logger.Errorln("缓存获取在线用户失败", err.Error())
			rsp.Fail("获取在线用户失败", ctx)
			return
		}

		// 转换为 LoginUser
		var loginUser login.LoginUser
		if err = json.Unmarshal([]byte(userStr), &loginUser); err != nil {
			gb.Logger.Errorln("将用户缓存信息转换为 LoginUser 对象失败", err.Error())
			rsp.Fail("获取在线用户失败", ctx)
			return
		}

		users = append(users, response.UserOnline{
			TokenId:   loginUser.Token,
			DeptName:  loginUser.DeptName,
			UserName:  loginUser.UserName,
			Ipaddr:    loginUser.IpAddress,
			Browser:   loginUser.Browser,
			Os:        loginUser.AgentOS,
			LoginTime: loginUser.LoginTime,
		})
	}

	// 在线用户列表
	users = util.NewList(users).Filter(func(uo response.UserOnline) bool {
		if userOnline.UserName != "" && uo.UserName != userOnline.UserName {
			return false
		}
		if userOnline.Ipaddr != "" && uo.Ipaddr != userOnline.Ipaddr {
			return false
		}
		return true
	}).Order(func(a, b response.UserOnline) int {
		return cmp.Compare(a.LoginTime, b.LoginTime)
	})

	// vpn 用户列表
	vpnUsers := util.Convert(openvpn.OpenvpnService.VpnStatus.OnlineUsers,
		func(u entity.SysLoginLog) response.UserOnline {
			return response.UserOnline{
				UserName:      u.UserName,
				Ipaddr:        u.Ipaddr,
				LoginLocation: "openvpn",
				LoginTime:     u.LoginTime.Time().UnixMilli(),
			}
		}).
		Order(func(a, b response.UserOnline) int {
			return cmp.Compare(a.LoginTime, b.LoginTime)
		})
	users = append(users, vpnUsers...)

	data := make(map[string]any)
	data["rows"] = users
	data["total"] = len(users)
	rsp.Context(ctx).Flat().OkWithData(data)
}

func (c *onlineController) ForceLogout(ctx *gin.Context) {
	tokenId := ctx.Param("tokenId")
	if tokenId != "" {
		gb.RedisProxy.Delete(cst.CACHE_LOGIN_TOKEN_KEY + tokenId)
		rsp.Ok(ctx)
		return
	}

	userName := ctx.Param("userName")
	if userName != "" {
		openvpn.OpenvpnService.KickOut(userName)
	}
	rsp.Ok(ctx)
}
