package login

import (
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"github.com/pkg/errors"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/request"
	"vpn-web.funcworks.net/service/login"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util/rsp"
)

var UserLogin = &loginController{}

type loginController struct {
	controller.BaseController
}

func (c *loginController) Login(ctx *gin.Context) {
	// 提取请求数据
	req, err := getLoginRequest(ctx)
	if err != nil {
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 校验验证码
	if !login.CaptchaService.Verify(req.Uuid, req.Code, true) {
		rsp.Fail("验证码不正确", ctx)
		return
	}

	// 用户登录
	token, err := login.LoginService.Login(req)
	if err != nil {
		gb.Logger.Errorf("用户[%s]登录失败: %s", req.Username, err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 返回 token
	data := make(map[string]interface{})
	data["token"] = token
	rsp.Context(ctx).Flat().OkWithData(data)
	gb.Logger.Debugf("用户[%s]登录成功", req.Username)
}

func getLoginRequest(ctx *gin.Context) (*request.LoginRequest, error) {
	req := &request.LoginRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		return nil, errors.Wrap(err, "用户登录json解析失败")
	}
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}
	if req.Uuid == "" || req.Code == "" {
		return nil, errors.New("验证码不能为空")
	}

	ua := useragent.New(ctx.GetHeader("User-Agent"))
	req.ClientIp = ctx.ClientIP()
	req.Browser, _ = ua.Browser()
	req.Os = ua.OS()

	return req, nil
}

func (c *loginController) Logout(ctx *gin.Context) {
	loginUser := c.GetLoginUser(ctx)

	if err := login.TokenService.DelLoginUser(loginUser.Token); err != nil {
		gb.Logger.Errorf("用户[%s]登出失败: %s", loginUser.User.UserName, err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	rsp.OkWithMessage("系统退出成功", ctx)
}

func (c *loginController) GetUserInfo(ctx *gin.Context) {
	loginUser := c.GetLoginUser(ctx)

	// user
	data := make(map[string]interface{})
	data["user"] = loginUser.User

	// role
	roles, err := system.RoleService.GetUserRolePerms(loginUser.User)
	if err != nil {
		gb.Logger.Errorf("获取%s角色权限失败: %s", loginUser.User.UserName, err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	rolePerms := make([]string, 0)
	for k := range roles {
		rolePerms = append(rolePerms, k)
	}
	data["roles"] = rolePerms

	// menu
	menuPerms := make([]string, 0)
	for k := range loginUser.Permissions {
		menuPerms = append(menuPerms, k)
	}
	data["permissions"] = menuPerms

	rsp.Context(ctx).Flat().OkWithData(data)
}

func (c *loginController) GetRouters(ctx *gin.Context) {
	loginUser := c.GetLoginUser(ctx)

	if menus, err := system.MenuService.GetUserMenuTree(loginUser); err != nil {
		gb.Logger.Errorf("获取%s菜单权限失败: %s", loginUser.User.UserName, err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.OkWithData(menus, ctx)
	}
}
