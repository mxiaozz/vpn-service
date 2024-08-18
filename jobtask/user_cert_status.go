package jobtask

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	model "vpn-web.funcworks.net/model/openvpn"
	service "vpn-web.funcworks.net/service/openvpn"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
)

func init() {
	gb.Sched.Registry("refreshUserCertStatus", refreshUserCertStatus)
}

// 更新用户证书有效期剩余天数
func refreshUserCertStatus(params []any, ctx context.Context) (any, error) {
	gb.Logger.Info("开始更新用户证书有效期剩余天数")

	users, err := system.UserService.GetAllUsers()
	if err != nil {
		return nil, errors.Wrap(err, "获取用户列表失败")
	}
	for _, user := range users {
		userCert, err := service.OpenvpnService.GetUserCert(user.UserName, false)
		if err != nil {
			return "", errors.Wrap(err, "获取用户证书失败")
		}
		err = updateUserCertValidDay(user.UserName, userCert)
		if err != nil {
			return "", errors.Wrap(err, "更新用户证书剩余天数失败")
		}
	}

	gb.Logger.Infof("完成 %d 个用户证书有效期更新", len(users))
	return "ok", nil
}

func updateUserCertValidDay(userName string, userCert *model.UserCert) error {
	validDay := "-"
	if !userCert.EndTime.IsZero() {
		days := util.DiffDays(userCert.EndTime, time.Now())
		validDay = fmt.Sprintf("%d天", days)
	}
	gb.Logger.Debugf("%s 证书有效期: %s", userName, validDay)

	user := &entity.SysUser{
		UserName: userName,
		ValidDay: validDay,
	}
	err := system.UserService.UpdateUserCertValidDay(user)
	if err != nil {
		return err
	}

	return nil
}
