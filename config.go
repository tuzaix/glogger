package glogger

/*
	通用日志切割代码，自动按照配置后缀自动切
*/

const (
	debugLevel = iota
	infoLevel
	noticeLevel
	warningLevel
	errorLevel
)

const (
	nocolor = 0
	red     = 30 + iota
	green
	yellow
	blue
	purple
	cyan
)

var (
	logLevel = map[int]string{
		debugLevel:   "DEBUG",
		infoLevel:    "INFO",
		noticeLevel:  "NOTICE",
		warningLevel: "WARNING",
		errorLevel:   "ERROR",
	}

	logColor = map[int]int{
		debugLevel:   cyan,
		infoLevel:    nocolor,
		noticeLevel:  green,
		warningLevel: yellow,
		errorLevel:   red,
	}

	LogStr2Int = map[string]int{
		"DEBUG":   debugLevel,
		"INFO":    infoLevel,
		"NOTICE":  noticeLevel,
		"WARNING": warningLevel,
		"ERROR":   errorLevel,
	}

	// 日志等级关系
	levelsMapping = map[string][]string{
		"DEBUG":   []string{"DEBUG", "NOTICE", "INFO", "WARNING", "ERROR"},
		"NOTICE":  []string{"NOTICE", "INFO", "WARNING", "ERROR"},
		"INFO":    []string{"INFO", "WARNING", "ERROR"},
		"WARNING": []string{"WARNING", "ERROR"},
		"ERROR":   []string{"ERROR"},
	}

	// 时间配置映射关系
	timeConfMapping = map[string]string{
		"%Y%m%d":       "20060102",
		"%Y%m%d%H":     "2006010215",
		"%Y%m%d%H%M":   "200601021504",
		"%Y%m%d%H%M%S": "20060102150405",

		"%Y-%m-%d":          "2006-01-02",
		"%Y-%m-%d-%H":       "2006-01-02-15",
		"%Y-%m-%d-%H-%M":    "2006-01-02-15-04",
		"%Y-%m-%d-%H-%M-%S": "2006-01-02-15-04-05",

		"%Y-%m-%d %H":       "2006-01-02 15",
		"%Y-%m-%d %H:%M":    "2006-01-02 15:04",
		"%Y-%m-%d %H:%M:%S": "2006-01-02 15:04:05",

		"%Y%m%d %H":       "20060102 15",
		"%Y%m%d %H:%M":    "20060102 15:04",
		"%Y%m%d %H:%M:%S": "20060102 15:04:05",
	}
)
