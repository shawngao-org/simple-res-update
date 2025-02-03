package logger

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-colorable"
	"time"
)

const (
	Green   = "\033[97;42m"
	White   = "\033[90;47m"
	Yellow  = "\033[97;43m"
	Red     = "\033[97;41m"
	Blue    = "\033[97;44m"
	Magenta = "\033[97;45m"
	Cyan    = "\033[97;46m"
	Reset   = "\033[0m"
)

func LogErr(format string, a ...any) {
	logTypeImpl("error", format, a...)
}

func LogWarn(format string, a ...any) {
	logTypeImpl("warn", format, a...)
}

func LogInfo(format string, a ...any) {
	logTypeImpl("info", format, a...)
}

func LogSuccess(format string, a ...any) {
	logTypeImpl("success", format, a...)
}

func LogDebug(format string, a ...any) {
	logTypeImpl("debug", format, a...)
}

func getNowTimeString() string {
	return time.Now().String()[0:19]
}

func logTypeImpl(typeStr string, format string, a ...any) {
	color := Blue
	tip := " I "
	switch typeStr {
	case "error":
		color = Red
		tip = " E "
	case "warn":
		color = Yellow
		tip = " W "
	case "success":
		color = Green
		tip = " U "
	case "debug":
		color = Magenta
		tip = " D "
	}
	logImpl(Reset+"["+getNowTimeString()+"] |"+color+tip+Reset+"|"+": "+format+"\n"+Reset, a...)
}

func logImpl(format string, a ...any) {
	stdOut := bufio.NewWriter(colorable.NewColorableStdout())
	_, err := fmt.Fprintf(stdOut, format, a...)
	if err != nil {
		return
	}
	_ = stdOut.Flush()
}
