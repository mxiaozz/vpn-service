package base

import (
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"vpn-web.funcworks.net/gb"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// 初始化Sqlite数据库
func initSqlite() {
	engine, err := xorm.NewEngine("sqlite3", gb.Config.Sqlite.Path)
	if err != nil {
		panic(err)
	}
	engine.SetLogger(&UlogXorm{SugaredLogger: gb.Logger})
	engine.ShowSQL(true)
	gb.DB = engine

	gb.Logger.Debug("initialized sqlite db")
}

type UlogXorm struct {
	*zap.SugaredLogger

	isShowSQL bool

	// zap:
	// DebugLevel  = -1
	// InfoLevel   = 0
	// WarnLevel   = 1
	// ErrorLevel  = 2
	// DPanicLevel = 3
	// PanicLevel  = 4
	// FatalLevel  = 5
}

func (ulog *UlogXorm) Level() log.LogLevel {
	// LOG_DEBUG   = 0
	// LOG_INFO    = 1
	// LOG_WARNING = 2
	// LOG_ERR     = 3
	// LOG_OFF     = 4
	// LOG_UNKNOWN = 5

	level := int(ulog.SugaredLogger.Level())
	switch level {
	case -1:
		return 0
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 3
	case 3:
		return 3
	case 4:
		return 3
	case 5:
		return 3
	default:
		return 1
	}
}

func (ulog *UlogXorm) SetLevel(l log.LogLevel) {

}

func (ulog *UlogXorm) ShowSQL(show ...bool) {
	ulog.isShowSQL = len(show) > 0 && show[0]
}

func (ulog *UlogXorm) IsShowSQL() bool {
	return ulog.isShowSQL
}
