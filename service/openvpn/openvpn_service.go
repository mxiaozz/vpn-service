package openvpn

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/buffer"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	vpn "vpn-web.funcworks.net/model/openvpn"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
)

var OpenvpnService = &openvpnService{
	VpnStatus: &vpn.OpenVpnStatus{
		Status:          "未知",
		LastUpdatedTime: model.DateTime(time.Now()),
	},
	mgmtUrl: gb.Viper.GetString("openvpn.mgmtUrl"),
}

type openvpnService struct {
	VpnStatus *vpn.OpenVpnStatus
	mgmtUrl   string
}

func (os *openvpnService) GetServerStatus() *vpn.OpenVpnStatus {
	return os.VpnStatus
}

func (os *openvpnService) GetRealtimeStatus() *vpn.OpenVpnStatus {
	status := &vpn.OpenVpnStatus{Status: "未知"}

	err := os.HandleServerStatus()
	if err == nil {
		status.Status = os.VpnStatus.Status
	}

	return status
}

// 获取 OpenVPN 进程运行状态
func (os *openvpnService) HandleServerStatus() error {
	data, err := util.HttpGet[string](os.mgmtUrl + "/serverStatus")
	if err != nil {
		return errors.Wrap(err, "http获取OpenVpn服务进程状态失败")
	}
	gb.Logger.Debugf("server status: %s", *data)

	switch *data {
	case "unknown":
		os.VpnStatus.Status = "未知"
	case "stopped":
		os.VpnStatus.Status = "停止"
	case "running":
		os.VpnStatus.Status = "正常"
	case "notinit":
		os.VpnStatus.Status = "未初始化"
	}

	return nil
}

func (os *openvpnService) ChangeServer(opt string) {
	var err error

	switch opt {
	case "start":
		_, err = util.HttpPost[string](os.mgmtUrl+"/serverStart", nil)
	case "stop":
		_, err = util.HttpPost[string](os.mgmtUrl+"/serverStop", nil)
	}
	if err != nil {
		gb.Logger.Errorf("changeServer: %s", err.Error())
	}

	// 暂停2s，等待 openvpn 状态稳定
	time.Sleep(2 * time.Second)
}

func (os *openvpnService) GetServerConfig() (string, error) {
	obj, err := util.HttpSend[string]("get", os.mgmtUrl+"/getConfig", nil)
	if err != nil {
		return "", errors.Wrap(err, "http获取OpenVpn服务进程状态失败")
	}
	if obj.Code != 0 && obj.Code != 2000 {
		return "", errors.Wrap(errors.New(obj.Msg), "http获取OpenVpn服务进程状态失败")
	}
	return obj.Data, nil
}

// 生成 OpenVPN 配置文件
func (os *openvpnService) GenerateConfig(param string) (string, error) {
	params, err := os.cmdParamPerse(param)
	if err != nil {
		return "", errors.Wrap(err, "生成OpenVPN配置shell指令解析失败")
	}
	gb.Logger.Infof("生成配置参数: %s", params)

	// 请求生成配置文件
	data := make(map[string]any, 0)
	data["params"] = params
	if _, err = util.HttpPost[string](os.mgmtUrl+"/genConfig", data); err != nil {
		return "", errors.Wrap(err, "http请求生成OpenVPN配置失败")
	}

	// 读取配置文件
	if cfg, err := util.HttpGet[string](os.mgmtUrl + "/getConfig"); err != nil {
		return "", errors.Wrap(err, "生成OpenVPN配置成功后读取配置失败")
	} else {
		return *cfg, nil
	}
}

// shell 命令行解析
func (os *openvpnService) cmdParamPerse(param string) ([]string, error) {
	params := make([]string, 0)

	// 前个字符
	preChar := rune(0)
	// 双引号中
	dQuoteFlag := false
	// 单引号中
	sQuoteFlag := false
	p := buffer.Buffer{}
	for _, c := range param {
		// 双引号
		if c == 34 {
			if sQuoteFlag {
				return nil, errors.New("参数中单引号包含双引号")
			}
			dQuoteFlag = !dQuoteFlag
			continue
		}
		// 单引号
		if c == 39 {
			if dQuoteFlag {
				return nil, errors.New("参数中双引号包含单引号")
			}
			sQuoteFlag = !sQuoteFlag
			continue
		}
		// 非空格为参数内容
		if c != 32 {
			p.AppendByte(byte(c))
			preChar = c
			continue
		}
		// 以下为空格处理逻辑
		// 空格在双引号或单引号中为参数内容
		if dQuoteFlag || sQuoteFlag {
			p.AppendByte(byte(c))
			continue
		}

		if preChar != 32 {
			params = append(params, p.String())
			p.Reset()
			preChar = c
			continue
		}
	}
	if p.Len() > 0 {
		params = append(params, p.String())
	}

	return params, nil
}

func (os *openvpnService) SaveConfig(cfgContent string) (string, error) {
	data := make(map[string]any, 0)
	data["content"] = cfgContent

	// 保存配置
	if _, err := util.HttpPost[string](os.mgmtUrl+"/saveConfig", data); err != nil {
		return "", errors.Wrap(err, "http保存OpenVPN配置失败")
	}

	// 读取配置
	if cfg, err := util.HttpGet[string](os.mgmtUrl + "/getConfig"); err != nil {
		return "", errors.Wrap(err, "保存OpenVPN配置成功后读取失败")
	} else {
		return *cfg, nil
	}
}

func (os *openvpnService) GetPKIStatus() (bool, error) {
	if status, err := util.HttpGet[bool](os.mgmtUrl + "/pkiStatus"); err != nil {
		return false, errors.Wrap(err, "http读取PKI状态失败")
	} else {
		return *status, nil
	}
}

func (os *openvpnService) InikPKI() error {
	gb.Logger.Info("初始化 OpenVpn PKI")
	if _, err := util.HttpSend[string]("post", os.mgmtUrl+"/initPKI", nil, func(c *http.Client, r *http.Request) {
		c.Timeout = 120 * time.Second
	}); err != nil {
		return errors.Wrap(err, "http初始化PKI失败")
	}
	return nil
}

func (os *openvpnService) ResetPKI() error {
	gb.Logger.Warn("重置 OpenVpn PKI")

	if _, err := util.HttpSend[string]("post", os.mgmtUrl+"/resetPKI", nil, func(c *http.Client, r *http.Request) {
		c.Timeout = 120 * time.Second
	}); err != nil {
		return errors.Wrap(err, "http重置PKI失败")
	}
	return nil
}

// 获取用户证书
func (os *openvpnService) GetUserCert(userName string, isWithCert bool) (*vpn.UserCert, error) {
	obj, err := util.HttpSend[vpn.UserCert]("get", os.mgmtUrl+"/getClientCert?name="+userName, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http获取用户证书失败")
	}

	// 2001：证书未注册
	if obj.Code == 2001 {
		return &vpn.UserCert{}, nil
	} else if obj.Code != 0 {
		return nil, errors.Wrap(errors.New(obj.Msg), "获取用户证书失败")
	}

	// 从客户端配置中解析出证书内容
	userCert := obj.Data
	isCertContent := false
	x509CertContent := buffer.Buffer{}
	reader := bufio.NewReader(strings.NewReader(userCert.Cert))
	for {
		lineStr, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "读取用户证书内容失败")
		}

		if strings.HasPrefix(lineStr, "<cert>") {
			isCertContent = true
			continue
		}
		if strings.HasPrefix(lineStr, "</cert>") {
			isCertContent = false
			break
		}
		if isCertContent {
			x509CertContent.AppendString(lineStr)
		}
	}

	// 证书读取
	block, _ := pem.Decode(x509CertContent.Bytes())
	if block == nil {
		return nil, errors.Wrap(err, "用户证书pem格式解码失败")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "用户证书x509解析失败")
	}
	userCert.Name = cert.Subject.CommonName
	userCert.BeginTime = model.DateTime(cert.NotBefore)
	userCert.EndTime = model.DateTime(cert.NotAfter)
	userCert.Durtion = fmt.Sprintf("%d 天", util.DiffDays(userCert.EndTime.Time(), time.Now()))
	if cert.NotAfter.After(time.Now()) {
		userCert.Status = "有效"
	} else {
		userCert.Status = "过期"
	}
	if !isWithCert {
		userCert.Cert = ""
	}

	gb.Logger.Infof("%s period: %s ~ %s, %s",
		userCert.Name,
		userCert.BeginTime.String(),
		userCert.EndTime.String(),
		userCert.Durtion)

	return &userCert, nil
}

func (os *openvpnService) GenerateUserCert(userName string) error {
	// 检查是否是系统用户
	if user, err := system.UserService.GetSysUser(userName, false); err != nil {
		return errors.Wrap(err, "生成用户证书查询用户失败")
	} else if user == nil {
		return errors.New("用户: " + userName + " 不存在，无法签发证书")
	}

	// 从库中读取有效期时长参数
	expire := 30
	if expireStr, err := system.ConfigService.GetConfigByKey("openvpn.cert.expire"); err != nil {
		return errors.Wrap(err, "读取用户证书有效期时长配置参数失败")
	} else if expireStr != "" {
		if expire, err = strconv.Atoi(expireStr); err != nil {
			return errors.Wrap(err, "生成用户证书有效期解析失败")
		}
	}

	gb.Logger.Infof("begin generate user: %s cert", userName)

	// 请求生成证书
	data := make(map[string]any, 0)
	data["name"] = userName
	data["expire"] = expire
	if _, err := util.HttpPost[string](os.mgmtUrl+"/genClientCert", data); err != nil {
		return errors.Wrap(err, "生成用户证书失败")
	}
	return nil
}

func (os *openvpnService) RevokeUserCert(userName string) error {
	gb.Logger.Infof("begin revoke user: %s cert", userName)
	data := make(map[string]any, 0)
	data["name"] = userName
	if _, err := util.HttpPost[string](os.mgmtUrl+"/revokeClientCert", data); err != nil {
		return errors.Wrap(err, "注销用户证书失败")
	}
	return nil
}

func (os *openvpnService) KickOut(userName string) error {
	gb.Logger.Infof("begin kick out user: %s", userName)
	data := make(map[string]any, 0)
	data["name"] = userName
	if obj, err := util.HttpPost[vpn.MgmtResponse](os.mgmtUrl+"/killClient", data); err != nil {
		return errors.Wrap(err, "强退VPN用户失败")
	} else {
		gb.Logger.Infof("kick out %s %s", userName, obj.Rsp)
		os.VpnStatus.OnlineUsers = slices.DeleteFunc(os.VpnStatus.OnlineUsers, func(u entity.SysLoginLog) bool {
			return u.UserName == userName
		})
	}
	return nil
}
