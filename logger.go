package dailyTask

/*
 * 用于打日志到生产机，与c++程序完全一样的日志格式
 *
 * 用法：
 * SetLogLevel(LogDebug) // 设置日志输出级别，默认为LogInfo
 * Info("now log level is %s", GetLogLevel()) // 获取当前日志级别
 *
 * Info("init daemon success, serverName %s, serverId %d", "myName", 123) // 打印日志
 * Err("something error")
 */

import (
	"github.com/rozen3/xLog"
)

func init() {
	LoggerInit("./dailyTask.log")
	SetLogLevel(ParseLogLevel("debug"))
}

const (
	LogMaxSize    = 1000
	LogMaxAge     = 0
	LogMaxBackups = 0
)

// 只能这样，如果再封装一层函数，打印出来的函数名就是Info,Err这种了
var Err = xLog.Error
var Warn = xLog.Warn
var Info = xLog.Info
var Debug = xLog.Debug

//var Err = log.Printf
//var Warn = log.Printf
//var Info = log.Printf
//var Debug = log.Printf

func ParseLogLevel(levelName string) int {
	switch levelName {
	case "error":
		return xLog.LevelError
	case "warn":
		return xLog.LevelWarn
	case "info":
		return xLog.LevelInfo
	case "debug":
		return xLog.LevelDebug
	default:
		return xLog.LevelInfo
	}
}

func LoggerInit(logFilePath string) {

	xLog.Init(logFilePath, LogMaxSize, LogMaxAge, LogMaxBackups, true, xLog.LevelInfo, false)

	SetLogLevel(xLog.LevelInfo)
}

func SetLogLevel(level int) {

	xLog.SetLevel(level)

	Info("set level to %s", xLog.LogLevelToString(level))
}

func GetLogLevel() int {
	return xLog.GetLevel()
}

// 默认情况下，调用日志函数都是往上查找caller并打印出caller的函数名
// 但是如果我们采用一些封装函数，会导致打印出来的函数名是封装函数，而不是更上一层的caller
// 这个时候就需要增加查找caller的深度了
//func ErrWithMoreDepth(depth int, args ...interface{}) {
//	//if myLogger.level < LogErr {
//	//	return
//	//}
//
//	Err(args)
//	//msg := fmt.Sprintf(args[0].(string), args[1:]...)
//
//	//myLogger.logger.Err(getHeader(LoggerCallDepth+depth) + msg)
//}
//
//func InfoWithMoreDepth(depth int, args ...interface{}) {
//	//if myLogger.level < LogErr {
//	//	return
//	//}
//
//	Info(args)
//	//msg := fmt.Sprintf(args[0].(string), args[1:]...)
//
//	//myLogger.logger.Info(getHeader(LoggerCallDepth+depth) + msg)
//}
//
//func getHeader(depth int) string {
//	pc, file, line, ok := runtime.Caller(depth)
//	if !ok {
//		file = "???"
//		line = 1
//	} else {
//		slash := strings.LastIndex(file, "/")
//		if slash >= 0 {
//			file = file[slash+1:]
//		}
//	}
//	funName := runtime.FuncForPC(pc).Name()
//	point := strings.LastIndex(funName, ".")
//	if point > 0 {
//		funName = string(funName[point+1:])
//	}
//
//	buff := fmt.Sprintf("%s:%d[%s] ", file, line, funName)
//
//	return buff
//}
