package system

import (
	"time"

	"github.com/reugn/go-quartz/quartz"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"xorm.io/builder"
)

var JobService = &jobService{}

type jobService struct {
}

func (js *jobService) GetAllJob() ([]entity.SysJob, error) {
	var jobs = []entity.SysJob{}
	err := gb.DB.Find(&jobs)
	return jobs, err
}

func (js *jobService) GetJobListPage(job *entity.SysJob, page *model.Page[entity.SysJob]) error {
	return gb.SelectPage(page, func(sql *builder.Builder) builder.Cond {
		sql.Select("*").From("sys_job").
			Where(builder.If(job.JobName != "", builder.Like{"job_name", job.JobName}).
				And(builder.If(job.JobGroup != "", builder.Eq{"job_group": job.JobGroup})).
				And(builder.If(job.Status != "", builder.Eq{"status": job.Status})).
				And(builder.If(job.InvokeTarget != "", builder.Like{"invoke_target", job.InvokeTarget})))
		return builder.Expr("job_id")
	})
}

func (js *jobService) GetJob(jobId int64) (*entity.SysJob, error) {
	var job entity.SysJob
	if exist, err := gb.DB.Where("job_id = ?", jobId).Get(&job); err != nil || !exist {
		return nil, err
	}

	if cronTrigger, err := quartz.NewCronTriggerWithLoc(job.CronExpression, time.Local); err == nil {
		if next, err := cronTrigger.NextFireTime(time.Now().UTC().UnixNano()); err == nil {
			job.NextValidTime = time.Unix(next/int64(time.Second), 0)
		}
	}

	return &job, nil
}

func (js *jobService) AddJob(job *entity.SysJob) error {
	if _, err := gb.DB.Insert(job); err != nil {
		return err
	}
	return gb.Sched.ScheduleJob(job)
}

func (js *jobService) UpdateJob(job *entity.SysJob) error {
	if job.Status == "1" {
		if err := gb.Sched.PauseJob(job); err != nil {
			return err
		}
	} else {
		if err := gb.Sched.ScheduleJob(job); err != nil {
			return err
		}
	}
	_, err := gb.DB.Where("job_id = ?", job.JobId).Update(job)
	return err
}

func (js *jobService) ChangeStatus(job *entity.SysJob) error {
	dbJob, err := js.GetJob(job.JobId)
	if err != nil {
		return err
	}

	if job.Status == "0" {
		if err := gb.Sched.ResumeJob(dbJob); err != nil {
			return err
		}
	} else {
		if err := gb.Sched.PauseJob(dbJob); err != nil {
			return err
		}
	}
	_, err = gb.DB.Table("sys_job").Cols("status").Where("job_id = ?", job.JobId).Update(job)
	return err
}

func (js *jobService) RunJob(job *entity.SysJob) error {
	if job, err := js.GetJob(job.JobId); err != nil {
		return err
	} else {
		return gb.Sched.RunJob(job)
	}
}

func (js *jobService) DeleteJobs(jobIds []int64) error {
	for _, id := range jobIds {
		if job, err := js.GetJob(id); err != nil {
			return err
		} else if job != nil {
			if err := gb.Sched.DeleteJob(job); err != nil {
				return err
			}
		}
	}
	_, err := gb.DB.Table("sys_job").In("job_id", jobIds).Delete()
	return err
}
