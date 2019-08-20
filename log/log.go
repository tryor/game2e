package log

import (
	"fmt"
	"runtime"
	"time"
)

func init() {
	SetLevel(INFO)
}

var currentLevel Level

func SetLevel(lv Level) {
	currentLevel = lv
}

func Error(v ...interface{}) {
	Log(ERROR, v...)
}

func Warn(v ...interface{}) {
	Log(WARN, v...)
}

func Info(v ...interface{}) {
	Log(INFO, v...)
}

func Debug(v ...interface{}) {
	Log(DEBUG, v...)
}

func Errorf(format string, params ...interface{}) {
	Log(ERROR, fmt.Sprintf(format, params...))
}

func Warnf(format string, params ...interface{}) {
	Log(WARN, fmt.Sprintf(format, params...))
}

func Infof(format string, params ...interface{}) {
	Log(INFO, fmt.Sprintf(format, params...))
}

func Debugf(format string, params ...interface{}) {
	Log(DEBUG, fmt.Sprintf(format, params...))
}

func Log(lv Level, v ...interface{}) {
	if lv >= currentLevel {
		pc, file, line, ok := runtime.Caller(2)
		funcName := ""
		if ok {
			f := runtime.FuncForPC(pc)
			//			file, line := f.FileLine(f.Entry())
			//			fmt.Println(f.Name(), f.Entry(), file, line)
			funcName = f.Name()
		}
		fmt.Printf("[%v] [%v] %v(%v).%v %v\n", time.Now().Format("2006-01-02 15:04:05.999"), LevelText[lv], file, line, funcName, fmt.Sprint(v...))
	}
}

type Level int

const (
	DEBUG Level = 0
	INFO  Level = 1
	WARN  Level = 2
	ERROR Level = 3
)

var LevelText = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}
