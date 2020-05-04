package beego

import (
		"github.com/astaxie/beego/logs"
		"strings"
)

var (
		aliLSLogger     *logs.BeeLogger
		fileLogger      *logs.BeeLogger
		connLogger      *logs.BeeLogger
		esLogger        *logs.BeeLogger
		mailLogger      *logs.BeeLogger
		slackLogger     *logs.BeeLogger
		consoleLogger   *logs.BeeLogger
		multiFileLogger *logs.BeeLogger
		loggerLevelMap  = map[string]int{
				"alert":    logs.LevelAlert,
				"critical": logs.LevelCritical,
				"error":    logs.LevelError,
				"warn":     logs.LevelWarn,
				"notice":   logs.LevelNotice,
				"info":     logs.LevelInfo,
				"debug":    logs.LevelDebug,
		}
)

// 获取当前日志适配名
func GetLoggerAdapter() string {
		return GetKvString("log_adapter", logs.AdapterConsole)
}

// 获取日志对象
func GetLogger(names ...string) *logs.BeeLogger {
		var name = "default"
		if len(names) > 0 {
				name = names[0]
		}
		name = strings.ToLower(name)
		if name != "default" && name == GetLoggerAdapter() {
				name = "default"
		}
		switch name {
		case "default":
				return logs.GetBeeLogger()
		case logs.AdapterFile:
				return GetFileLogger()
		case logs.AdapterConn:
				return GetConnLogger()
		case logs.AdapterEs:
				return GetEsLogger()
		case logs.AdapterMail:
				return GetMailLogger()
		case logs.AdapterConsole:
				return GetConsoleLogger()
		case logs.AdapterSlack:
				return GetSlackLogger()
		case logs.AdapterMultiFile:
				return GetMultiFileLogger()
		case logs.AdapterAliLS:
				return GetAliLsLogger()
		}
		return logs.GetBeeLogger()
}

func GetAliLsLogger() *logs.BeeLogger {
		if aliLSLogger == nil {
				BootLoggerAdapter(logs.AdapterAliLS)
		}
		return aliLSLogger
}

func GetMultiFileLogger() *logs.BeeLogger {
		if multiFileLogger == nil {
				BootLoggerAdapter(logs.AdapterMultiFile)
		}
		return multiFileLogger
}

func GetSlackLogger() *logs.BeeLogger {
		if slackLogger == nil {
				BootLoggerAdapter(logs.AdapterSlack)
		}
		return slackLogger
}

func GetConsoleLogger() *logs.BeeLogger {
		if consoleLogger == nil {
				BootLoggerAdapter(logs.AdapterConsole)
		}
		return consoleLogger
}

func GetMailLogger() *logs.BeeLogger {
		if mailLogger == nil {
				BootLoggerAdapter(logs.AdapterMail)
		}
		return mailLogger
}

func GetEsLogger() *logs.BeeLogger {
		if esLogger == nil {
				BootLoggerAdapter(logs.AdapterEs)
		}
		return esLogger
}

func GetConnLogger() *logs.BeeLogger {
		if connLogger == nil {
				BootLoggerAdapter(logs.AdapterConn)
		}
		return connLogger
}

func GetFileLogger() *logs.BeeLogger {
		if fileLogger == nil {
				BootLoggerAdapter(logs.AdapterFile)
		}
		return fileLogger
}

// 初始化日志logger 适配器
func BootLoggerAdapter(adapters ...string) {
		var (
				isDef   = false
				logger  *logs.BeeLogger
				scope   string
				adapter string
		)
		if len(adapters) > 0 {
				adapter = adapters[0]
		}
		// 获取默认适配器
		if adapter == "" {
				adapter = GetLoggerAdapter()
		}
		// 获取日志对象
		if GetLoggerAdapter() != adapter {
				scope = adapter
				logger = logs.NewLogger()
		} else {
				isDef = true
				logger = logs.GetBeeLogger()
		}
		// 相关配置
		config := GetLoggerProperties(scope)
		_ = logger.SetLogger(adapter, config)
		logger.SetLevel(GetLoggerLevel())
		async := GetKvInt("log_async", 1)
		depth := GetKvInt("log_func_call_depth", 0)
		// 异步写入
		if async >= 1 {
				if async == 1 {
						logger.Async()
				} else {
						logger.Async(int64(async))
				}
		}
		// 日志追加调用层级
		if depth >= 2 {
				logger.EnableFuncCallDepth(true)
				if async > 2 {
						logger.SetLogFuncCallDepth(depth)
				}
		}
		switch adapter {
		case "default":
		case logs.AdapterFile:
				if !isDef {
						fileLogger = logger
				}
		case logs.AdapterConn:
				if !isDef {
						connLogger = logger
				}
		case logs.AdapterEs:
				if !isDef {
						esLogger = logger
				}
		case logs.AdapterMail:
				if !isDef {
						mailLogger = logger
				}
		case logs.AdapterConsole:
				if !isDef {
						consoleLogger = logger
				}
		case logs.AdapterSlack:
				if !isDef {
						slackLogger = logger
				}
		case logs.AdapterMultiFile:
				if !isDef {
						multiFileLogger = logger
				}
		case logs.AdapterAliLS:
				if !isDef {
						aliLSLogger = logger
				}
		}
}

// 获取日志配置参数
func GetLoggerProperties(scopes ...string) string {
		var scope string
		if len(scopes) > 0 {
				scope = scopes[0]
		}
		if scope == "" || scope == "default" {
				return GetParsedKvStr("log_properties", "{}")
		}
		return GetParsedKvStr("log_"+scope+"_properties", "{}")
}

// 日志监听等级
func GetLoggerLevel() int {
		level := GetKvString("log_level", "info")
		return GetLoggerLevelInt(level)
}

// 获取日志等级map
func GetLoggerLevelMap() map[string]int {
		return loggerLevelMap
}

// 通过等级名获取等级值
func GetLoggerLevelInt(level string) int {
		level = strings.ToLower(level)
		if v, ok := loggerLevelMap[level]; ok {
				return v
		}
		return 0
}
