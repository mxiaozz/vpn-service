package openvpn

import (
	"time"

	"vpn-web.funcworks.net/model/entity"
)

type OpenVpnStatus struct {
	Status          string               `json:"status"`
	StartTime       time.Time            `json:"startTime"`
	LastUpdatedTime time.Time            `json:"lastUpdatedTime"`
	Duration        string               `json:"duration"`
	OnlineUsers     []entity.SysLoginLog `json:"onlineUsers"`
}
