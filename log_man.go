package glogger

var (
	GlobalLoggers map[string]*Logger
)

func InitLogger(confFile string) {
	GlobalLoggers = NewLoggers(confFile)
}

func GetLogger(logName string) (logger *Logger) {
	// 获取logger，根据不同的名字
	var ok bool
	if logger, ok = GlobalLoggers[logName]; ok {
		return
	}
	return nil
}