package routers

import (
		_ "github.com/WebGameLinux/cms/middlewares"
		"github.com/WebGameLinux/cms/routers/api"
		"github.com/WebGameLinux/cms/routers/web"
)

func init() {
		api.RegisterApi("/v1")
		web.RegisterWeb("/")
}
