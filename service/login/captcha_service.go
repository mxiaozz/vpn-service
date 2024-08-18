package login

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

var CaptchaService = &base64Captcha.Captcha{
	Driver: &base64Captcha.DriverMath{
		Height:          38,
		Width:           115,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine,
		Fonts:           []string{"Times New Roman"},
	},
	Store: base64Captcha.NewMemoryStore(100, time.Duration(2*60*time.Second)),
}
