package login

import (
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/model/login"
	"vpn-web.funcworks.net/model/request"
	"vpn-web.funcworks.net/service/system"
)

var LoginService = &loginService{}

type loginService struct {
}

// 用户登录
func (us *loginService) Login(req request.LoginRequest) (string, error) {
	gb.Logger.Infof("用户登录：%s", req.Username)

	// 从数据库读取用户
	sysUser, err := us.getUserByName(req.Username)
	if err != nil {
		return "", err
	}

	// 检查登录错误次数
	count, err := us.checkLoginedErrCount(sysUser.UserName)
	if err != nil {
		return "", err
	}

	// 验证密码
	if err = us.checkPassword(sysUser.UserName, sysUser.Password, req.Password, count); err != nil {
		return "", errors.New("用户名或密码错误")
	}
	// 以防密码泄露
	sysUser.Password = ""

	// 加载登录用户权限
	loginUser, err := us.newLoginUser(sysUser)
	if err != nil {
		return "", errors.Wrap(err, "加载登录用户权限")
	}
	// 访问信息
	loginUser.IpAddress = req.ClientIp
	loginUser.Browser = req.Browser
	loginUser.AgentOS = req.Os

	// 生成登录 token，并缓存用户登录信息
	if jwtToken, err := TokenService.CreateToken(loginUser); err != nil {
		return "", errors.Wrap(err, "创建 token")
	} else {
		us.updateLoginInfo(sysUser.UserId, req.ClientIp)
		return jwtToken, nil
	}
}

func (us *loginService) getUserByName(userName string) (entity.SysUser, error) {
	user, err := system.UserService.GetSysUser(userName, true)
	if err != nil {
		return user, err
	}
	if user.Status == cst.USER_STATUS_DISABLE {
		return user, errors.New("用户已停用")
	}
	if user.DelFlag == cst.USER_STATUS_DELETED {
		return user, errors.New("用户已被删除")
	}
	return user, nil
}

func (us *loginService) checkLoginedErrCount(userName string) (int, error) {
	count, err := gb.RedisProxy.GetInt(us.getCacheErrKey(userName))
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return 0, err
	}

	if count >= gb.Config.Login.MaxRetryCount {
		gb.Logger.With().Warnf("登录用户：%s 密码错误次数超过限制.", userName)
		return 0, errors.New("密码错误次数超过限制")
	}

	return count, nil
}

func (us *loginService) updateLoginedErrCount(count int, userName string) error {
	if err := gb.RedisProxy.SetEx(
		us.getCacheErrKey(userName),
		count,
		gb.Config.Login.LockTime*60); err != nil {
		return err
	}
	return nil
}

func (us *loginService) checkPassword(userName, userPassword, password string, count int) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		us.updateLoginedErrCount(count+1, userName)
		return err
	}
	if count > 0 {
		gb.RedisProxy.Delete(us.getCacheErrKey(userName))
	}
	return nil
}

func (us *loginService) newLoginUser(sysUser entity.SysUser) (login.LoginUser, error) {
	var loginUser login.LoginUser

	perms, err := system.MenuService.GetMenuPermission(sysUser.UserId)
	if err != nil {
		return loginUser, err
	}

	loginUser.UserId = sysUser.UserId
	loginUser.UserName = sysUser.UserName
	loginUser.DeptId = sysUser.DeptId
	loginUser.DeptName = sysUser.Dept.DeptName
	loginUser.Permissions = perms
	return loginUser, nil
}

func (us *loginService) getCacheErrKey(userName string) string {
	return cst.CACHE_PWD_ERR_CNT_KEY + userName
}

func (us *loginService) updateLoginInfo(userId int64, clientIp string) {
	user := &entity.SysUser{
		UserId:    userId,
		LoginIp:   clientIp,
		LoginDate: model.DateTime(time.Now()),
	}
	if err := system.UserService.UpdateUserLoginInfo(user); err != nil {
		gb.Logger.Error("更新登录信息失败", err)
	}
}

func (us *loginService) Unlock(userName string) error {
	return gb.RedisProxy.Delete(us.getCacheErrKey(userName))
}
