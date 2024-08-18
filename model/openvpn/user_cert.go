package openvpn

import "time"

type UserCert struct {
	Name      string    `json:"name"`
	Cert      string    `json:"cert"`
	BeginTime time.Time `json:"beginTime"`
	EndTime   time.Time `json:"endTime"`
	Durtion   string    `json:"durtion"`
	Status    string    `json:"status"`
}
