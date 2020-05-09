package api

import (
		. "github.com/WebGameLinux/cms/controllers"
		"github.com/WebGameLinux/cms/middlewares"
		. "github.com/astaxie/beego"
)

// 自动注册

// 版本好接口
func RegisterApi(prefix string) {
		//	Router("/", &controllers.MainController{})
		// NSRouter("/login", new(controllers.MainController), "*:get"),
		ns := NewNamespace(prefix+"/user",
				NSRouter("/:id", GetUserController(), "get:GetUserById"),
		)
		// 无需登陆
		ns2 := NewNamespace(prefix,
				NSRouter("/login", GetLoginController(), "post:Login"),
				NSRouter("/register", GetUserController(), "post:Register"),
		)

		AddNamespace(ns)
		AddNamespace(ns2)
		ResisterApiMiddleware(prefix, "/user/*")
}

func ResisterApiMiddleware(prefix string, p ...string) {
		var uri string
		if len(p) != 0 {
				uri = prefix + p[0]
		} else {
				uri = prefix
		}
		InsertFilter(uri, BeforeExec, middlewares.Auth())
}
