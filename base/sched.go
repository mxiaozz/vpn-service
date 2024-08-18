package base

import (
	"github.com/reugn/go-quartz/logger"
	"go.uber.org/zap"
	"vpn-web.funcworks.net/gb"
	"vpn-web.funcworks.net/schedule"
)

func initSched() {
	logger.SetDefault(&UlogQuartz{SugaredLogger: gb.Logger})
	sched := schedule.NewSchedManager()
	sched.Init()
	gb.Sched = sched
}

type UlogQuartz struct {
	*zap.SugaredLogger

	// zap:
	// DebugLevel  = -1
	// InfoLevel   = 0
	// WarnLevel   = 1
	// ErrorLevel  = 2
	// DPanicLevel = 3
	// PanicLevel  = 4
	// FatalLevel  = 5
}

func (ulog *UlogQuartz) Trace(msg any) {

}

func (ulog *UlogQuartz) Tracef(format string, args ...any) {

}

func (ulog *UlogQuartz) Debug(msg any) {
	ulog.SugaredLogger.Debug(msg)
}

func (ulog *UlogQuartz) Info(msg any) {
	ulog.SugaredLogger.Info(msg)
}

func (ulog *UlogQuartz) Warn(msg any) {
	ulog.SugaredLogger.Warn(msg)
}

func (ulog *UlogQuartz) Error(msg any) {
	ulog.SugaredLogger.Error(msg)
}

func (ulog *UlogQuartz) Enabled(level logger.Level) bool {
	// LevelTrace  = -8
	// LevelDebug  = -4
	// LevelInfo   = 0
	// LevelWarn   = 4
	// LevelError  = 8
	// LevelOff    = 12

	level1 := int(ulog.SugaredLogger.Level())
	level2 := int(level)

	// debug
	if level1 == -1 {
		return level2 >= -4
	}
	// info
	if level1 == 0 {
		return level2 >= 0
	}
	//warn
	if level1 == 1 {
		return level2 >= 4
	}
	//error
	if level1 == 2 {
		return level2 >= 8
	}
	// panic | fatal
	if level1 >= 3 {
		return level2 >= 8
	}

	return false
}
