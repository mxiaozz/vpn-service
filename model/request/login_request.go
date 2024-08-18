package request

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Uuid     string `json:"uuid"`

	ClientIp string `json:"-"`
	Browser  string `json:"-"`
	Os       string `json:"-"`
}
