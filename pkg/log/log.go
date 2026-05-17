package log

import (
	"os"
	"sync"
	"time"

	ll "github.com/moorada/log"
)

var pathDir = "logs"
var pathFile = ""
var consoleMu sync.Mutex
var consoleEnabled bool

func ActiveConsoleLog() (err error) {
	consoleMu.Lock()
	defer consoleMu.Unlock()
	if consoleEnabled {
		return nil
	}

	config := ll.FormatConfigBasic
	config.Format = "{time} {level:color}{level:name}{reset} {message}"
	err = ll.AddOutput("", ll.INFO, config, false)
	if err != nil {
		return err
	}
	consoleEnabled = true
	return nil
}

func RemoveConsoleLog() {
	consoleMu.Lock()
	defer consoleMu.Unlock()
	ll.RemoveOutput("")
	consoleEnabled = false
}

func RemoveDebugLog() {
	ll.RemoveOutput(pathFile)
}

func ActiveDebugLog() (err error) {

	pathFile = pathDir + "/" + time.Now().Format("2January2006-15:04:05") + ".log"
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		_ = os.MkdirAll(pathDir, os.ModePerm)
	}

	config := ll.FormatConfigBasic
	config.Format = "{datetime} {level:color}{level:name}{reset} {message}"
	err = ll.AddOutput(pathFile, ll.DEBUG, config, true)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	ll.CloseOutputs()
}

func Debug(format string, args ...interface{}) {
	ll.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	ll.Info(format, args...)
}

func Important(format string, args ...interface{}) {
	ll.Important(format, args...)
}

func Warning(format string, args ...interface{}) {
	ll.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	ll.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	ll.Fatal(format, args...)
}
