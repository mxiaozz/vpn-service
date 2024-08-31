package system

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
	"vpn-web.funcworks.net/util/rsp"
)

var ProfileController = &profileController{}

type profileController struct {
	controller.BaseController
}

func (c *profileController) GetOwnerInfo(ctx *gin.Context) {
	loginUser := c.GetLoginUser(ctx)

	data := make(map[string]interface{})
	user, err := system.UserService.GetSysUserById(loginUser.UserId, true)
	if err != nil {
		gb.Logger.Errorln("获取用户信息失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	user.Password = ""
	data["user"] = user

	roles, err := system.RoleService.GetUserRoles(loginUser.UserId)
	if err != nil {
		gb.Logger.Errorln("获取用户角色失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	rlist := util.Convert(roles, func(r entity.SysRole) string { return r.RoleName })
	data["roleGroup"] = strings.Join(rlist, ",")

	posts, err := system.PostService.GetUserPostList(loginUser.UserId)
	if err != nil {
		gb.Logger.Errorln("获取用户岗位失败", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}
	plist := util.Convert(posts, func(p entity.SysPost) string { return p.PostName })
	data["postGroup"] = strings.Join(plist, ",")

	rsp.Context(ctx).Flat().OkWithData(data)
}

func (c *profileController) UpdateOwnerInfo(ctx *gin.Context) {
	var user entity.SysUser
	if err := ctx.ShouldBind(&user); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	user.UserId = loginUser.UserId

	if err := system.UserService.UpdateOwnerInfo(user); err != nil {
		gb.Logger.Errorln("修改个人信息失败", err.Error())
		rsp.Fail("修改个人信息异常，请联系管理员", ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *profileController) UpdateOwnerPassword(ctx *gin.Context) {
	var paramMap = make(map[string]string)
	if err := ctx.ShouldBind(&paramMap); err != nil {
		gb.Logger.Errorf(err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	newPassword := paramMap["newPassword"]
	oldPassword := paramMap["oldPassword"]

	if err := system.UserService.UpdateOwnerPassword(loginUser.UserId, newPassword, oldPassword); err != nil {
		gb.Logger.Errorln("修改个人密码失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Ok(ctx)
	}
}

func (c *profileController) UpdateOwnerAvatar(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		gb.Logger.Errorln("文件上传", err)
		rsp.Fail("文件上传失败", ctx)
		return
	}
	defer file.Close()

	// 检查文件后缀
	ext, err := c.checkAvatarFileType(header.Filename)
	if err != nil {
		gb.Logger.Errorln("头像文件", err.Error())
		rsp.Fail(err.Error(), ctx)
		return
	}

	// 检查文件大小
	if header.Size > 5*1024*1024 {
		gb.Logger.Errorln("头像文件超过5M")
		rsp.Fail("文件大小不能超过5M", ctx)
		return
	}

	loginUser := c.GetLoginUser(ctx)
	if url, err := system.UserService.UpdateOwnerAvatar(loginUser, file, ext); err != nil {
		gb.Logger.Errorln("修改个人头像失败", err.Error())
		rsp.Fail(err.Error(), ctx)
	} else {
		rsp.Context(ctx).Flat().OkWithData(gin.H{"imgUrl": url})
	}
}

// 检查文件格式
func (c *profileController) checkAvatarFileType(fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return "", errors.New("文件格式不正确")
	}
	fileTypeList := []string{".png", ".jpg", ".jpeg", ".gif"}
	if ext, exist := util.NewList(fileTypeList).First(func(t string) bool { return strings.EqualFold(ext, t) }); !exist {
		return "", errors.New("文件格式需为: png/jpg/jpeg/gif")
	} else {
		return ext, nil
	}
}
