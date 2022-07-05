package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	workingDir = "/"
	lock = sync.Mutex{}
	loggerMap = make(map[string]*log.Logger)
)
type MContextHook struct{}

func init() {
	wd, err := os.Getwd()
	if err == nil {
		workingDir = filepath.ToSlash(wd) + "/"
	}
}



func (hook *MContextHook) Levels() []log.Level {
	return []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel}
}

func (hook *MContextHook) getCallerInfo() (string, string, int) {
	var (
		shortPath string
		funcName  string
	)

	for i := 3; i < 15; i++ {
		pc, fullPath, line, ok := runtime.Caller(i)
		if !ok {
			continue
		} else {
			shortPath = fullPath
			funcName = runtime.FuncForPC(pc).Name()
			index := strings.LastIndex(funcName, ".")
			if index > 0 {
				funcName = funcName[index+1:]
			}
			if !strings.Contains(strings.ToLower(fullPath), "github.com/sirupsen/logrus") {
				return shortPath, funcName, line
				break
			}
		}
	}
	return "", "", 0
}

func (hook *MContextHook) Fire(entry *log.Entry) error {
	shortPath, funcName, callLine := hook.getCallerInfo()
	if shortPath != "" && callLine != 0 {
		entry.Data["caller"] = fmt.Sprintf("[%s:%s:%d]", shortPath, funcName, callLine)
	}
	return nil
}


func GetLogger(filePath string) *log.Logger {
	lock.Lock()
	defer lock.Unlock()
	pathParts := strings.Split(filePath, "/")
	fileName := pathParts[len(pathParts)-1]
	if logger, ok := loggerMap[fileName]; ok {
		return logger
	} else {
		logger := log.New()
		logger.SetLevel(log.WarnLevel)
		logger.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			DisableColors:   true,
			FullTimestamp:   true,
		})
		logger.AddHook(&MContextHook{})
		loggerMap[fileName] = logger
		return logger
	}
}
