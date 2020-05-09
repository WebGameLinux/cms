package bootstarp

import (
		_ "github.com/WebGameLinux/cms/routers" // 注册router
)

func init() {
		// 初始化日志输出收集
		bootLogger()
		// 先加载models
		bootModels()
		// 后初始化数据库
		bootDatabase()
		// 初始化 worker
		bootWorker()
}
