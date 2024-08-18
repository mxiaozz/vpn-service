package jobtask

import (
	"context"

	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/service/system"
)

func init() {
	gb.Sched.Registry("cleanJobLog", cleanJobLog)
}

// 保留最近N天日志记录
func cleanJobLog(params []any, ctx context.Context) (any, error) {
	days := 3
	if len(params) > 0 {
		if v, ok := params[0].(int); ok {
			days = v
		}
	}
	system.JobLogService.CleanJobLog(days)
	return nil, nil
}
