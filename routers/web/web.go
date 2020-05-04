package web

import (
		"github.com/WebGameLinux/cms/controllers"
		"github.com/astaxie/beego"
)

func init() {
		// beego.Router("/", &controllers.MainController{})
}

func RegisterWeb(prefix string) {
		beego.Router(prefix, &controllers.MainController{})
}
