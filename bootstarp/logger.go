package bootstarp

import (
		utils "github.com/WebGameLinux/cms/utils/beego"
)

// 初始化默认 logger 适配器
func bootLogger() {
		utils.BootLoggerAdapter()
}
