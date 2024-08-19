package openvpn

import (
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
)

type OpenVpnStatus struct {
	Status          string               `json:"status"`
	StartTime       model.DateTime       `json:"startTime"`
	LastUpdatedTime model.DateTime       `json:"lastUpdatedTime"`
	Duration        string               `json:"duration"`
	OnlineUsers     []entity.SysLoginLog `json:"onlineUsers"`
}
