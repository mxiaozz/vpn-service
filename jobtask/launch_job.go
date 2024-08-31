package jobtask

import (
	"strconv"

	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/service/system"
)

// 程序启动时，加载所有定时任务
func LoadSchedJobs() {
	// 注册任务
	registryJob()

	// 加载任务
	if jobs, err := system.JobService.GetAllJob(); err != nil {
		gb.Logger.Errorln("任务调度加载所有任务失败", err)
	} else {
		for _, job := range jobs {
			if job.Status == "1" {
				continue
			}

			schedJob := gb.SchedJob{
				JobId:          strconv.FormatInt(job.JobId, 10),
				JobName:        job.JobName,
				JobGroup:       job.JobGroup,
				InvokeTarget:   job.InvokeTarget,
				CronExpression: job.CronExpression,
			}

			if err := gb.Sched.ScheduleJob(schedJob); err != nil {
				gb.Logger.Errorf("%s 任务调度失败: %s", job.JobName, err.Error())
			}
		}
	}
	gb.Logger.Info("调度任务启动完成")
}
