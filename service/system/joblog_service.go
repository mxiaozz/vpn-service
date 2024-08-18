package system

import (
	"time"

	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
	"xorm.io/xorm"
)

var JobLogService = &jobLogService{}

type jobLogService struct {
}

func (jls *jobLogService) GetJobLogListPage(jobLog *entity.SysJobLog, page *model.Page[entity.SysJobLog]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_job_log").
			Where(builder.If(jobLog.JobName != "", builder.Like{"job_name", jobLog.JobName}).
				And(builder.If(jobLog.JobGroup != "", builder.Eq{"job_group": jobLog.JobGroup})).
				And(builder.If(jobLog.Status != "", builder.Eq{"status": jobLog.Status})).
				And(builder.If(jobLog.InvokeTarget != "", builder.Like{"invoke_target", jobLog.InvokeTarget})).
				And(builder.If(func() bool { return jobLog.Params["beginTime"] != nil }(),
					builder.Gte{"create_time": jobLog.Params["beginTime"]})).
				And(builder.If(func() bool { return jobLog.Params["endTime"] != nil }(),
					builder.Lte{"create_time": jobLog.Params["endTime"]})))
		return builder.Expr("job_log_id desc")
	})
}

func (jls *jobLogService) GetJobLog(jobLogId int64) (*entity.SysJobLog, error) {
	var job entity.SysJobLog
	if exist, err := gb.DB.Where("job_log_id = ?", jobLogId).Get(&job); err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}
	return &job, nil
}

func (jls *jobLogService) AddJobLog(jobLog *entity.SysJobLog) error {
	_, err := gb.DB.Insert(jobLog)
	return err
}

func (jls *jobLogService) DeleteJobLogs(jobLogIds []int64) error {
	_, err := gb.DB.Table("sys_job_log").In("job_log_id", jobLogIds).Delete()
	return err
}

func (jls *jobLogService) CleanJobLogs() error {
	return gb.Tx(func(dbSession *xorm.Session) error {
		if _, err := dbSession.Exec("delete from sys_job_log"); err != nil {
			return err
		}
		_, err := dbSession.Exec("delete from sqlite_sequence where name = 'sys_job_log'")
		return err
	})
}

// 清除N天前日志
func (jls *jobLogService) CleanJobLog(days int) {
	lastDate := time.Now().Add(time.Duration(-days) * 24 * time.Hour).Format("2006-01-02")
	if _, err := gb.DB.Table("sys_job_log").Where("create_time < ?", lastDate).Delete(); err != nil {
		gb.Logger.Error(err)
	}
}
