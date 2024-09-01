package openvpn

import (
	"vpn-web.funcworks.net/controller"
	"vpn-web.funcworks.net/router/wraper"
)

func InitOpenVpnRouter(pvt wraper.RouterWraper) {
	vpn := wraper.ExtModule("VPN管理")
	pvt = wraper.RouterWraper{RouterGroup: pvt.Group("/openvpn")}
	{
		pvt.GET("/getstatus", controller.Openvpn.GetStatus, vpn.Ext())
		pvt.GET("/getRealtimeStatus", controller.Openvpn.GetRealtimeStatus, vpn.Ext())
		pvt.GET("/getServerConfig", controller.Openvpn.GetServerConfig, vpn.Ext())
		pvt.POST("/generateConfig", controller.Openvpn.GenerateConfig, vpn.Ext("openvpn:config:generate"))
		pvt.PUT("/saveConfig", controller.Openvpn.SaveConfig, vpn.Ext("openvpn:config:save"))
		pvt.GET("/getPKIStatus", controller.Openvpn.GetPKIStatus, vpn.Ext())
		pvt.POST("/inikPKI", controller.Openvpn.InikPKI, vpn.Ext("openvpn:pki:init"))
		pvt.PUT("/resetPKI", controller.Openvpn.ResetPKI, vpn.Ext("openvpn:pki:reset"))
		pvt.PUT("/optServer", controller.Openvpn.OptServer, vpn.Ext("openvpn:server:opt"))
		pvt.GET("/getUserCert", controller.Openvpn.GetUserCert, vpn.Ext())
		pvt.POST("/generateUserCert", controller.Openvpn.GenerateUserCert, vpn.Ext("openvpn:cert:generate"))
		pvt.DELETE("/revokeUserCert", controller.Openvpn.RevokeUserCert, vpn.Ext("openvpn:cert:revoke"))
		pvt.POST("/downUserCert", controller.Openvpn.DownUserCert, vpn.Ext("openvpn:cert:download"))
	}
}
