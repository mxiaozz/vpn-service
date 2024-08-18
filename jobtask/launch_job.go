package jobtask

import (
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/service/system"
)

// 程序启动时，加载所有定时任务
func LoadSchedJobs() {
	if jobs, err := system.JobService.GetAllJob(); err != nil {
		gb.Logger.Errorln("任务调度加载所有任务失败", err)
	} else {
		for _, job := range jobs {
			if job.Status == "1" {
				continue
			}
			if err := gb.Sched.ScheduleJob(&job); err != nil {
				gb.Logger.Errorf("%s 任务调度失败: %s", job.JobName, err.Error())
			}
		}
	}
	gb.Logger.Info("调度任务启动完成")
}
