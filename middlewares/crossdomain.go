package middlewares

import (
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/plugins/cors"
)

func init() {
		crossDomain()
}

// 参数
func Options() *cors.Options {
		return &cors.Options{
				AllowHeaders:     getAllowHeaders([]string{"*"}),
				AllowAllOrigins:  getAllowAll(false),
				AllowOrigins:     getAllowOrigins([]string{"localhost"}),
				AllowMethods:     getAllowMethods([]string{"PUT", "PATCH", "POST", "GET", "DELETE"}),
				ExposeHeaders:    getExposeHeaders([]string{"Content-Length"}),
				AllowCredentials: getAllowCredentials(true),
		}
}

// 跨域
func crossDomain() {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(Options()))
}

// 是否运行所有
func getAllowAll(defaultValue ...bool) bool {
		return utils.GetKvBool("cross_allow_all", defaultValue...)
}

// 获取
func getAllowHeaders(defaultValue ...[]string) []string {
		return utils.GetKvStrArr("cross_allow_headers", defaultValue...)
}

func getAllowOrigins(defaultValue ...[]string) []string {
		return utils.GetKvStrArr("cross_domains", defaultValue...)
}

func getAllowMethods(defaultValue ...[]string) []string {
		return utils.GetKvStrArr("cross_allow_methods", defaultValue...)
}

func getAllowCredentials(defaultValue ...bool) bool {
		return utils.GetKvBool("cross_allow_credentials", defaultValue...)
}

func getExposeHeaders(defaultValue ...[]string) []string {
		return utils.GetKvStrArr("cross_expose_headers", defaultValue...)
}
