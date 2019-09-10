package glogger

/*
	通用日志切割代码，自动按照配置后缀自动切
*/

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

type LoggerConf struct {
	LogName       string `toml:"log_name"`
	LogLevel      string `toml:"log_level"`
	LogDir        string `toml:"log_dir"`
	LogFile       string `toml:"log_file"`
	LogReserve    int    `toml:"log_reserve"`
	LogDataFormat string `toml:"log_date_format"`
	LogColor      int    `toml:"log_color"`
	LogConsole    int    `toml:"log_console"`
}

type LoggerConfs struct {
	LConfs map[string]LoggerConf `toml:"loggers"`
}

func NewLoggerConf(logConf string) map[string]LoggerConf {
	var (
		data  []byte
		err   error
		confs LoggerConfs
	)
	data, err = ioutil.ReadFile(logConf)
	if err != nil {
		panic("加载日志配置文件失败" + err.Error())
	}

	if _, err = toml.Decode(string(data), &confs); err != nil {
		panic("toml解析配置失败" + err.Error())
	}
	return confs.LConfs
}

type Logger struct {
	sync.Mutex
	isConsole      bool
	isColorfull    bool
	reserveCounter int
	timeFormat     string
	fileName       string
	fileWriter     io.WriteCloser

	Levels map[int]bool
	mtime  string
}

func NewLoggers(confFile string) (loggers map[string]*Logger, loggerConfs map[string]LoggerConf) {
	loggerConfs = NewLoggerConf(confFile)
	loggers = make(map[string]*Logger)
	logLoggers := make([]*Logger, 0)
	for _, conf := range loggerConfs {
		tmp_logger := NewLogger(conf.LogLevel, conf.LogDir, conf.LogFile, conf.LogReserve, conf.LogDataFormat, conf.LogConsole, conf.LogColor)
		loggers[conf.LogName] = tmp_logger
		logLoggers = append(logLoggers, tmp_logger)
	}
	go logSplitChecker(logLoggers)
	go logCounterChecker(logLoggers)
	return
}

func NewLogger(level string, dir string, file string, reserve int, suffix string, console int, color int) *Logger {
	logFilePath := fmt.Sprintf("%s/%s", dir, file)
	fileWriter := getFileWriter(logFilePath)
	boolConsole := false
	if console == 1 {
		boolConsole = true
	}
	boolColorfull := false
	if color == 1 {
		boolColorfull = true
	}

	var (
		timeFormat string
		ok         bool
	)

	if timeFormat, ok = timeConfMapping[suffix]; !ok {
		panic("日志后缀错误")
		return nil
	}

	loggerHandle := &Logger{
		isConsole:      boolConsole,
		isColorfull:    boolColorfull,
		reserveCounter: reserve,
		timeFormat:     timeFormat,
		fileName:       logFilePath,
		fileWriter:     fileWriter,
	}

	// 获取不同优先级的日志级别
	if levels, ok := levelsMapping[level]; ok {
		loggerHandle.setLevel(levels)
	}
	return loggerHandle
}

// 设置日志级别
func (l *Logger) setLevel(levels []string) {
	l.Levels = make(map[int]bool)
	for _, v := range levels {
		if sv, ok := LogStr2Int[v]; ok {
			l.Levels[sv] = true
		}
	}
}

func (l *Logger) write(level int, format string, content ...interface{}) {
	filename, fc, line := getDetail() // 获取文件的信息
	if len(l.Levels) > 0 {
		if _, ok := l.Levels[level]; !ok {
			return
		}
	}
	now := time.Now()
	var s string
	if format == "" {
		s = renderColor(fmt.Sprintf("%s ◊ %s ◊ %s ◊ %s:%d ◊ %s\n", now.Format("2006/01/02 15:04:05"), logLevel[level], filename, fc, line, fmt.Sprint(content...)), logColor[level], l.isColorfull)
	} else {
		s = renderColor(fmt.Sprintf("%s ◊ %s ◊ %s ◊ %s:%d ◊ %s\n", now.Format("2006/01/02 15:04:05"), logLevel[level], filename, fc, line, fmt.Sprintf(format, content...)), logColor[level], l.isColorfull)
	}

	l.Lock()
	defer l.Unlock()
	l.fileWriter.Write([]byte(s))
	if l.isConsole {
		fmt.Print(s)
	}
}

func (l *Logger) Info(content ...interface{}) {
	l.write(infoLevel, "", content...)
}

func (l *Logger) Infof(format string, content ...interface{}) {
	l.write(infoLevel, format, content...)
}

func (l *Logger) Warning(content ...interface{}) {
	l.write(warningLevel, "", content...)
}

func (l *Logger) Warningf(format string, content ...interface{}) {
	l.write(warningLevel, format, content...)
}

func (l *Logger) Notice(content ...interface{}) {
	l.write(noticeLevel, "", content...)
}

func (l *Logger) Noticef(format string, content ...interface{}) {
	l.write(noticeLevel, format, content...)
}

func (l *Logger) Debug(content ...interface{}) {
	l.write(debugLevel, "", content...)
}

func (l *Logger) Debugf(format string, content ...interface{}) {
	l.write(debugLevel, format, content...)
}

func (l *Logger) Error(content ...interface{}) {
	l.write(errorLevel, "", content...)
}

func (l *Logger) Errorf(format string, content ...interface{}) {
	l.write(errorLevel, format, content...)
}

func renderColor(s string, color int, isColorfull bool) string {
	if isColorfull {
		return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	} else {
		return s
	}
}
