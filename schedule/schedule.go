package schedule

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	quartzJob "github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"vpn-web.funcworks.net/cst"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/model"
	"vpn-web.funcworks.net/model/entity"
	"vpn-web.funcworks.net/service/system"
	"vpn-web.funcworks.net/util"
)

type schedManager struct {
	sched       quartz.Scheduler
	workFuncMap map[string]gb.WorkFuncion
}

func NewSchedManager() *schedManager {
	return &schedManager{
		sched:       quartz.NewStdScheduler(),
		workFuncMap: make(map[string]gb.WorkFuncion),
	}
}

func (sm *schedManager) Init() {
	sm.sched.Start(context.Background())
}

func (sm *schedManager) Registry(methodName string, workFunc gb.WorkFuncion) error {
	if _, ok := sm.workFuncMap[methodName]; ok {
		return errors.New(methodName + " job already exists")
	}
	sm.workFuncMap[methodName] = workFunc
	return nil
}

func (sm *schedManager) ScheduleJob(job gb.SchedJob) error {
	jobKey := sm.getJobKey(job)
	if schedJob, _ := sm.sched.GetScheduledJob(jobKey); schedJob != nil {
		sm.sched.DeleteJob(jobKey)
	}

	var trigger quartz.Trigger
	if job.CronExpression != "" {
		cronTrigger, err := quartz.NewCronTriggerWithLoc(job.CronExpression, time.Local)
		if err != nil {
			return err
		}
		trigger = cronTrigger
	} else if job.Interval > 0 {
		trigger = quartz.NewSimpleTrigger(time.Duration(job.Interval) * time.Second)
	} else {
		return errors.New("job.cronExpression and job.interval is empty")
	}

	funcJob, err := sm.newFunctionJob(job)
	if err != nil {
		return err
	}
	jobDetail := quartz.NewJobDetail(funcJob, jobKey)

	gb.Logger.Infoln("增加调度任务:", job.JobName)
	return sm.sched.ScheduleJob(jobDetail, trigger)
}

func (sm *schedManager) getJobKey(job gb.SchedJob) *quartz.JobKey {
	return quartz.NewJobKey(fmt.Sprintf("task-%s-%s", job.JobGroup, job.JobId))
}

func (sm *schedManager) newFunctionJob(job gb.SchedJob) (*quartzJob.FunctionJob[any], error) {
	methodName, params := sm.parseInvokeMethod(job.InvokeTarget)
	if methodName == "" {
		return nil, errors.New("methodName: " + methodName + " in job.invokeTarget is empty")
	}

	var workFunc gb.WorkFuncion
	if wfunc, ok := sm.workFuncMap[methodName]; !ok {
		return nil, errors.New(methodName + " function not found")
	} else {
		workFunc = wfunc
	}

	return quartzJob.NewFunctionJob(func(ctx context.Context) (any, error) {
		startTime := time.Now()

		//exec
		rstValue, rstErr := workFunc(params, ctx)

		// log
		stopTime := time.Now()
		sysJogLob := entity.SysJobLog{
			JobName:      job.JobName,
			JobGroup:     job.JobGroup,
			InvokeTarget: job.InvokeTarget,
			CreateTime:   model.DateTimeNow(),
		}
		if rstErr != nil {
			sysJogLob.Status = cst.SYS_FAIL
		} else {
			sysJogLob.Status = cst.SYS_SUCCESS
		}
		cost := stopTime.Sub(startTime).Milliseconds()
		sysJogLob.JobMessage = fmt.Sprintf("%s 总共耗时：%d 毫秒", job.JobName, cost)
		if err := system.JobLogService.AddJobLog(sysJogLob); err != nil {
			gb.Logger.Errorln("调度任务记录日志失败", job.JobName, err)
		}

		return rstValue, errors.Wrap(rstErr, "调度任务执行失败: "+job.JobName)
	}), nil
}

func (sm *schedManager) parseInvokeMethod(method string) (string, []any) {
	if method == "" {
		return "", nil
	}
	method = strings.Trim(method, " ")
	params := make([]any, 0)

	start := strings.Index(method, "(")
	if start == -1 {
		return method, params
	}

	end := strings.Index(method, ")")
	if end == -1 {
		return method, params
	}

	args := strings.Split(method[start+1:end], ",")
	if len(args) == 0 {
		return method, params
	}

	return method[:start], util.Convert(args, func(v string) any {
		v = strings.Trim(v, " ")
		if strings.HasPrefix(v, "\"") || strings.HasPrefix(v, "'") {
			return v[1 : len(v)-1]
		}
		if strings.EqualFold(v, "true") {
			return true
		}
		if strings.EqualFold(v, "false") {
			return false
		}
		if strings.HasSuffix(v, "L") || strings.HasSuffix(v, "l") {
			l, _ := strconv.ParseInt(v[:len(v)-1], 10, 64)
			return l
		}
		if strings.HasSuffix(v, "D") || strings.HasSuffix(v, "d") {
			d, _ := strconv.ParseFloat(v[:len(v)-1], 64)
			return d
		}
		if strings.HasSuffix(v, "F") || strings.HasSuffix(v, "f") {
			d, _ := strconv.ParseFloat(v[:len(v)-1], 64)
			return float32(d)
		}
		i, _ := strconv.ParseInt(v, 10, 32)
		return int(i)
	})
}

func (sm *schedManager) PauseJob(job gb.SchedJob) error {
	jobKey := sm.getJobKey(job)
	if schedJob, _ := sm.sched.GetScheduledJob(jobKey); schedJob == nil {
		return nil
	} else if schedJob.JobDetail().Options().Suspended {
		return nil
	}
	gb.Logger.Infoln("暂定调度任务:" + job.JobName)
	return sm.sched.PauseJob(jobKey)
}

func (sm *schedManager) ResumeJob(job gb.SchedJob) error {
	jobKey := sm.getJobKey(job)
	if schedJob, _ := sm.sched.GetScheduledJob(jobKey); schedJob == nil {
		return sm.ScheduleJob(job)
	} else if schedJob.JobDetail().Options().Suspended {
		gb.Logger.Infoln("恢复调度任务:" + job.JobName)
		return sm.sched.ResumeJob(jobKey)
	}
	return nil
}

func (sm *schedManager) DeleteJob(job gb.SchedJob) error {
	jobKey := sm.getJobKey(job)
	if schedJob, _ := sm.sched.GetScheduledJob(jobKey); schedJob == nil {
		return nil
	} else {
		gb.Logger.Infoln("删除调度任务:" + job.JobName)
		return sm.sched.DeleteJob(jobKey)
	}
}

func (sm *schedManager) RunJob(job gb.SchedJob) error {
	if workFunc, err := sm.newFunctionJob(job); err != nil {
		return err
	} else {
		gb.Logger.Infoln("执行一次调度任务:" + job.JobName)
		return workFunc.Execute(context.Background())
	}
}
