package login

import (
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/model/login"
	"vpn-web.funcworks.net/model/request"
	"vpn-web.funcworks.net/service/system"
)

var LoginService = &loginService{}

type loginService struct {
}

// 用户登录
func (us *loginService) Login(req *request.LoginRequest) (string, error) {
	gb.Logger.Infof("用户登录：%s", req.Username)

	// 从数据库读取用户
	sysUser, err := us.getUserByName(req.Username)
	if err != nil {
		return "", err
	}

	// 检查登录错误次数
	count, err := us.checkLoginedErrCount(sysUser)
	if err != nil {
		return "", err
	}

	// 验证密码
	if err = us.checkPassword(req.Password, count, sysUser); err != nil {
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
	loginUser.Os = req.Os

	// 生成登录 token，并缓存用户登录信息
	if jwtToken, err := TokenService.CreateToken(loginUser); err != nil {
		return "", errors.Wrap(err, "创建 token")
	} else {
		us.updateLoginInfo(sysUser.UserId, req.ClientIp)
		return jwtToken, nil
	}
}

func (us *loginService) getUserByName(userName string) (*entity.SysUser, error) {
	user, err := system.UserService.GetSysUser(userName, true)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status == cst.USER_STATUS_DISABLE {
		return nil, errors.New("用户已停用")
	}
	if user.DelFlag == cst.USER_STATUS_DELETED {
		return nil, errors.New("用户已被删除")
	}
	return user, nil
}

func (us *loginService) checkLoginedErrCount(user *entity.SysUser) (int, error) {
	count, err := gb.RedisProxy.GetInt(us.getCacheErrKey(user.UserName))
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return 0, err
	}

	if count >= gb.Config.Login.MaxRetryCount {
		gb.Logger.With().Warnf("登录用户：%s 密码错误次数超过限制.", user.UserName)
		return 0, errors.New("密码错误次数超过限制")
	}

	return count, nil
}

func (us *loginService) updateLoginedErrCount(count int, user *entity.SysUser) error {
	if err := gb.RedisProxy.SetEx(
		us.getCacheErrKey(user.UserName),
		count,
		gb.Config.Login.LockTime*60); err != nil {
		return err
	}
	return nil
}

func (us *loginService) checkPassword(password string, count int, user *entity.SysUser) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		us.updateLoginedErrCount(count+1, user)
		return err
	}
	if count > 0 {
		gb.RedisProxy.Delete(us.getCacheErrKey(user.UserName))
	}
	return nil
}

func (us *loginService) newLoginUser(sysUser *entity.SysUser) (*login.LoginUser, error) {
	perms, err := system.MenuService.GetMenuPermission(sysUser)
	if err != nil {
		return nil, err
	}

	loginUser := &login.LoginUser{
		UserId:      sysUser.UserId,
		DeptId:      sysUser.DeptId,
		User:        sysUser,
		Permissions: perms,
	}
	return loginUser, nil
}

func (us *loginService) getCacheErrKey(userName string) string {
	return cst.CACHE_PWD_ERR_CNT_KEY + userName
}

func (us *loginService) updateLoginInfo(userId int64, clientIp string) {
	user := &entity.SysUser{
		UserId:    userId,
		LoginIp:   clientIp,
		LoginDate: time.Now(),
	}
	if err := system.UserService.UpdateUserLoginInfo(user); err != nil {
		gb.Logger.Error("更新登录信息失败", err)
	}
}

func (us *loginService) Unlock(userName string) error {
	return gb.RedisProxy.Delete(us.getCacheErrKey(userName))
}
