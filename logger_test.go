package glogger_test

import (
	. "github.com/tuzaix/glogger"
	"testing"
)

func TestNewLoggerConf(t *testing.T) {

	InitLogger("log.toml")

    GetLogger().Info("默认logger.........")
    GLogger("logger2").Info("非默认logger.........")

	select {}

}
