package jobtask

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model/entity"
	model "vpn-web.funcworks.net/model/openvpn"
	service "vpn-web.funcworks.net/service/openvpn"
	"vpn-web.funcworks.net/util"
)

var vpnStausSchedule = &openvpnSchedule{
	vpnStatus: service.OpenvpnService.VpnStatus,
	mgmtUrl:   gb.Viper.GetString("openvpn.mgmtUrl"),
}

type openvpnSchedule struct {
	vpnStatus *model.OpenVpnStatus
	mgmtUrl   string
}

// 定时刷新 OpenVPN 服务状态
func (os *openvpnSchedule) refreshVpnStatus(params []any, ctx context.Context) (any, error) {
	gb.Logger.Info("begin refresh openvpn status")

	os.vpnStatus.Status = "未知"
	os.vpnStatus.OnlineUsers = nil

	// 获取服务状态
	if err := service.OpenvpnService.HandleServerStatus(); err != nil {
		return nil, err
	}
	// 获取服务启动时间
	if err := os.handleState(); err != nil {
		return nil, err
	}
	// 获取服务状态 和 当前在线用户
	if err := os.handleOnlineUsers(); err != nil {
		return nil, err
	}

	gb.Logger.Infof("end refresh openvpn status\n%+v", os.vpnStatus)
	return "ok", nil
}

// 获取 OpenVPN 服务启动时间
func (os *openvpnSchedule) handleState() error {
	obj, err := util.HttpGet[model.MgmtResponse](os.mgmtUrl + "/state")
	if err != nil {
		return errors.Wrap(err, "http获取服务启动时间失败")
	}
	data := obj.Rsp
	gb.Logger.Debugf("state response: \n%s", data)

	// 状态信息
	// 1721888674,CONNECTED,SUCCESS,10.254.250.1,,,,\r\nEND\r\n
	reader := bufio.NewReader(strings.NewReader(data))
	stateStr, err := reader.ReadString('\n')
	if err != io.EOF && err != nil {
		return errors.Wrap(err, "解析服务启动内容失败")
	}
	if stateStr == "" || stateStr == "END" || strings.HasPrefix(stateStr, "ERROR") {
		return nil
	}
	dateStr := strings.Split(stateStr, ",")[0]
	timestamp, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "解析服务启动时间失败: %s", dateStr)
	}
	os.vpnStatus.StartTime = time.Unix(timestamp, 0)

	return nil
}

// 获取服务状态 和 当前在线用户
func (os *openvpnSchedule) handleOnlineUsers() error {
	obj, err := util.HttpGet[model.MgmtResponse](os.mgmtUrl + "/status")
	if err != nil {
		return errors.Wrap(err, "http获取当前在线用户失败")
	}
	data := obj.Rsp
	gb.Logger.Debugf("state response: \n%s", data)

	statusStr := ""
	statusFlag := false
	userFlag := false
	routerFlag := false
	reader := bufio.NewReader(strings.NewReader(data))
	for {
		statusStr, err = reader.ReadString('\n')
		if err == io.EOF && statusStr == "" {
			break
		}
		if err != nil {
			return errors.Wrap(err, "解析当前在线用户内容失败")
		}
		statusStr = strings.TrimRightFunc(statusStr, func(r rune) bool {
			return r == '\r' || r == '\n'
		})

		// 读取结束 或 遇到错误
		if statusStr == "END" || strings.HasPrefix(statusStr, "ERROR") {
			break
		}

		if strings.HasPrefix(statusStr, "OpenVPN CLIENT LIST") {
			statusFlag = true
			userFlag = false
			routerFlag = false
			continue
		}
		if strings.HasPrefix(statusStr, "Common Name") {
			userFlag = true
			statusFlag = false
			routerFlag = false
			continue
		}
		if strings.HasPrefix(statusStr, "Virtual Address") {
			routerFlag = true
			userFlag = false
			statusFlag = false
			continue
		}

		// 解析服务状态
		if statusFlag {
			array := strings.Split(statusStr, ",")
			if len(array) != 2 {
				statusFlag = false
				continue
			}

			dateStr := array[1]
			date, err := time.Parse("2006-01-02 15:04:05", dateStr)
			if err != nil {
				return errors.Wrapf(err, "解析服务器当前时间失败: %s", dateStr)
			}
			os.vpnStatus.Duration = util.TimeDistance(date, os.vpnStatus.StartTime)
			statusFlag = false
		}

		// 解析用户列表
		if userFlag {
			array := strings.Split(statusStr, ",")
			if len(array) != 5 {
				continue
			}

			sizeSend, err := strconv.ParseInt(array[2], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "解析用户发送数据量失败: %s", array[2])
			}
			sizeReceive, err := strconv.ParseInt(array[3], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "解析用户接受数据量失败: %s", array[3])
			}
			loginTime, err := time.Parse("2006-01-02 15:04:05", array[4])
			if err != nil {
				return errors.Wrapf(err, "解析用户登录时间失败: %s", array[4])
			}
			user := entity.SysLoginLog{
				UserName:  array[0],
				Ipaddr:    strings.Split(array[1], ":")[0],
				Browser:   util.HumanByteSize(sizeSend),
				Os:        util.HumanByteSize(sizeReceive),
				LoginTime: loginTime,
				Msg:       util.TimeDistance(time.Now(), loginTime),
			}
			if os.vpnStatus.OnlineUsers == nil {
				os.vpnStatus.OnlineUsers = make([]entity.SysLoginLog, 0)
			}
			os.vpnStatus.OnlineUsers = append(os.vpnStatus.OnlineUsers, user)
		}

		// 解析路由
		if routerFlag {
			array := strings.Split(statusStr, ",")
			if len(array) != 4 {
				continue
			}
			if os.vpnStatus.OnlineUsers == nil {
				continue
			}

			for _, u := range os.vpnStatus.OnlineUsers {
				if u.UserName == array[1] {
					u.LoginLocation = array[0]
					// 免除编辑器告警：unused write to field LoginLocation
					_ = u.LoginLocation
					break
				}
			}
		}
	}

	return nil
}
