package openvpn

import (
	"vpn-web.funcworks.net/model"
)

type UserCert struct {
	Name      string         `json:"name"`
	Cert      string         `json:"cert"`
	BeginTime model.DateTime `json:"beginTime"`
	EndTime   model.DateTime `json:"endTime"`
	Durtion   string         `json:"durtion"`
	Status    string         `json:"status"`
}
