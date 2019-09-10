package glogger_test

import (
	. "glogger"
	"testing"
	"time"
)

func routine(nl map[string]*Logger) {
	ticker := time.NewTicker(100 * time.Microsecond) // 文件切割
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lg1, _ := nl["logger1"]
			lg1.Error("fdsafdfadfadfa")
			lg2, _ := nl["logger2"]
			lg2.Info("fdsafdfadfadfa")
		}
	}
}

func TestNewLoggerConf(t *testing.T) {

	InitLogger("log.toml")

	go routine(GlobalLoggers)
	go routine(GlobalLoggers)
	go routine(GlobalLoggers)
	go routine(GlobalLoggers)
	go routine(GlobalLoggers)

	select {}

}
