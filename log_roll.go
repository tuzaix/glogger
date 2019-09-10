package glogger

/*
	日志切割模块
*/

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	mtime string
)

// 创建文件
func getFileWriter(fileName string) io.WriteCloser {
	fileWriter, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return fileWriter
}

func getDetail() (string, string, int) {
	// 获取调用函数的数据
	pc, file, line, _ := runtime.Caller(3)
	fc := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fc, ".")
	filename := path.Base(file)
	return fmt.Sprintf("%s/%s", parts[0], filename), parts[1], line
}

type FileStat struct {
	name   string
	fmtime int64
}
type FSTS []FileStat

func (fs FSTS) Len() int           { return len(fs) }
func (fs FSTS) Swap(i, j int)      { fs[i], fs[j] = fs[j], fs[i] }
func (fs FSTS) Less(i, j int) bool { return fs[i].fmtime > fs[j].fmtime }

// 滚动切割文件
func rollingLogFile(toFileName string, logger *Logger) {
	logger.Lock()
	defer logger.Unlock()

	logger.fileWriter.Close()
	logger.fileWriter = nil

	err := os.Rename(logger.fileName, toFileName)
	if err != nil {
		//panic(err)
	}
	fileWriter, err := os.OpenFile(logger.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		//panic(err)
	}
	logger.fileWriter = fileWriter
}

func logSplitChecker(loggers []*Logger) {
	// log split checker
	ticker := time.NewTicker(10 * time.Second) // 文件切割
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, lg := range loggers {
				if lg.mtime == "" {
					lg.mtime = time.Now().Format(lg.timeFormat)
					continue
				}
				currentTime := time.Now().Format(lg.timeFormat)
				if currentTime != lg.mtime {
					toFileName := fmt.Sprintf("%s.%s", lg.fileName, lg.mtime)
					rollingLogFile(toFileName, lg)
					lg.mtime = currentTime
				}
			}
		}
	}
}

func logCounterChecker(loggers []*Logger) {
	// 保持日志个数
	ticker := time.NewTicker(600 * time.Second) // 一分钟检查一次日志个数
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, lg := range loggers {
				dirname := filepath.Dir(lg.fileName)
				basename := filepath.Base(lg.fileName)
				logLists := make(FSTS, 0)
				filepath.Walk(fmt.Sprintf("%s/", dirname), func(path string, f os.FileInfo, err error) error {
					if strings.HasPrefix(f.Name(), fmt.Sprintf("%s.", basename)) {
						fs := FileStat{
							name:   f.Name(),
							fmtime: f.ModTime().Unix(),
						}
						logLists = append(logLists, fs)
					}
					return nil
				})
				sort.Sort(logLists)

				if len(logLists) > lg.reserveCounter {
					removes := logLists[lg.reserveCounter:]
					for _, fname := range removes {
						rmname := fmt.Sprintf("%s/%s", dirname, fname.name)
						os.Remove(rmname)
					}
				}
			}
		}
	}
}
