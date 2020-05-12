package api

import (
		. "github.com/WebGameLinux/cms/controllers"
		"github.com/WebGameLinux/cms/middlewares"
		. "github.com/astaxie/beego"
)

// 自动注册

// 版本好接口
func RegisterApi(prefix string) {
		// 需要登陆 用户信息接口
		AddNamespace(NewNamespace(prefix+"/user",
				NSRouter("/:id", GetUserController(), "get:GetUserById"),
		))
		// 用户无需登陆
		AddNamespace(NewNamespace(prefix,
				NSRouter("/login", GetLoginController(), "post:Login"),
				NSRouter("/register", GetUserController(), "post:Register"),
		))
		// 附件接口
		AddNamespace(
				NewNamespace("/attachment",
						NSRouter("/:id", GetAttachmentController(), "get:GetById"),
						NSRouter("/lists", GetAttachmentController(), "get:Lists"),
						NSRouter("/upload", GetAttachmentController(), "post:Upload"),
						NSRouter("/uploads", GetAttachmentController(), "post:Uploads"),
				),
		)
		// 用户登陆验证
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
