package api

import (
		"github.com/WebGameLinux/cms/controllers"
		. "github.com/astaxie/beego"
)

// 自动注册

// 版本好接口
func RegisterApi(prefix string) {
		//	Router("/", &controllers.MainController{})
		// NSRouter("/login", new(controllers.MainController), "*:get"),
		c := new(controllers.UserController)
		ns := NewNamespace(prefix+"/user",
				NSRouter("/:id", c, "get:GetUserById"),
		)

		AddNamespace(ns)
}
