package system

import (
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var OperLogService = &operLogService{}

type operLogService struct {
}

func (os *operLogService) GetOperLogListPage(operLog entity.SysOperLog, page *model.Page[entity.SysOperLog]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_oper_log").
			Where(builder.If(operLog.OperIp != "", builder.Like{"oper_ip", operLog.OperIp}).
				And(builder.If(operLog.Title != "", builder.Like{"title", operLog.Title})).
				And(builder.If(operLog.BusinessType > -1, builder.Eq{"business_type": operLog.BusinessType})).
				And(builder.If(operLog.Status > -1, builder.Eq{"status": operLog.Status})).
				And(builder.If(operLog.OperName != "", builder.Like{"oper_name", operLog.OperName})).
				And(builder.If(func() bool { return operLog.Params["beginTime"] != nil }(),
					builder.Gte{"oper_time": operLog.Params["beginTime"]})).
				And(builder.If(func() bool { return operLog.Params["endTime"] != nil }(),
					builder.Lte{"oper_time": operLog.Params["endTime"]})))
		return builder.Expr("oper_id desc")
	})
}

func (os *operLogService) AddOperLog(operLog entity.SysOperLog) error {
	_, err := gb.DB.Insert(operLog)
	return err
}

func (os *operLogService) DeleteOperLogs(operLogIds []int64) error {
	_, err := gb.DB.Table("sys_oper_log").In("oper_id", operLogIds).Delete()
	return err
}

func (os *operLogService) CleanOperLogs() error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		if _, err := dbSession.Exec("delete from sys_oper_log"); err != nil {
			return err
		}
		_, err := dbSession.Exec("delete from sqlite_sequence where name = 'sys_oper_log'")
		return err
	})
}
