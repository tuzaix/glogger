package glogger

import (
    "sort"
    "strconv"
)

var (
    // 全局logger对象
	GlobalLoggers map[string]*Logger
    // 全局logger配置
    GlobalLoggerConfs map[string]LoggerConf
    // 默认logger
    defaultLoggerName string
)

func InitLogger(confFile string) {
	GlobalLoggers, GlobalLoggerConfs = NewLoggers(confFile)
    initDefaultLoggerName()
}

func initDefaultLoggerName() {
    // 设置第一个key为默认logger
    if len(GlobalLoggerConfs) == 0 {
        panic("没有配置文件...")
    }
    keys := make([]int, len(GlobalLoggerConfs))
    j := 0
    for k := range GlobalLoggerConfs {
        i_k, err := strconv.Atoi(k)
        if err != nil {
            continue
        }
        keys[j] = i_k
        j++
    }
    sort.Ints(keys)
    defaultIndex := strconv.Itoa(keys[0])
    if defaultLoggerConf, ok := GlobalLoggerConfs[defaultIndex]; ok {
        defaultLoggerName = defaultLoggerConf.LogName
    }
}

func GetLogger(logNames ...string) (logger *Logger) {
	// 获取logger，根据不同的名字
    var logName string
    if len(logNames) == 0 {
        // 使用默认第一个logger 
        logName = defaultLoggerName
    } else {
        // 无论传几个，只有第一个LogName生效
        logName = logNames[0]
    }
	var ok bool
	if logger, ok = GlobalLoggers[logName]; ok {
		return
	}
	return nil
}

