package monitor

import (
	"strconv"
	"strings"

	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/util"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var LoginLogService = &loginLogService{}

type loginLogService struct {
}

func (ls *loginLogService) GetLoginLogListPage(loginLog entity.SysLoginLog, page *model.Page[entity.SysLoginLog]) error {
	if err := gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_logininfor").
			Where(builder.If(loginLog.UserName != "", builder.Eq{"user_name": loginLog.UserName}).
				And(builder.If(loginLog.Status != "", builder.Eq{"status": loginLog.Status})).
				And(builder.If(loginLog.Ipaddr != "", builder.Like{"ipaddr", loginLog.Ipaddr})).
				And(builder.If(func() bool { return loginLog.Params["beginTime"] != nil }(),
					builder.Gte{"login_time": loginLog.Params["beginTime"]})).
				And(builder.If(func() bool { return loginLog.Params["endTime"] != nil }(),
					builder.Lte{"login_time": loginLog.Params["endTime"]})))
		return builder.Expr("info_id desc")
	}); err != nil {
		return err
	}

	for i, log := range page.Rows {
		if strings.HasPrefix(log.Browser, "[send]") {
			sendValueStr := strings.Trim(log.Browser[6:], " ")
			if sendValueStr != "" {
				if sendBytes, err := strconv.ParseInt(sendValueStr, 10, 64); err == nil {
					page.Rows[i].Browser = "[send] " + util.HumanByteSize(sendBytes)
				}
			}
		}
		if strings.HasPrefix(log.Os, "[recv]") {
			recvValueStr := strings.Trim(log.Os[6:], " ")
			if recvValueStr != "" {
				if recvBytes, err := strconv.ParseInt(recvValueStr, 10, 64); err == nil {
					page.Rows[i].Os = "[recv] " + util.HumanByteSize(recvBytes)
				}
			}
		}
	}

	return nil
}

func (ls *loginLogService) AddLoginLog(loginLog entity.SysLoginLog) error {
	_, err := gb.DB.Insert(loginLog)
	return err
}

func (ls *loginLogService) DeleteLoginLogs(loginLogIds []int64) error {
	_, err := gb.DB.Table("sys_logininfor").In("info_id", loginLogIds).Delete()
	return err
}

func (ls *loginLogService) CleanLoginLogs() error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		if _, err := dbSession.Exec("delete from sys_logininfor"); err != nil {
			return err
		}
		_, err := dbSession.Exec("delete from sqlite_sequence where name = 'sys_logininfor'")
		return err
	})
}
