package bootstarp

import (
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/models/captcha"
		"github.com/astaxie/beego/orm"
)

// 注册model
func bootModels() {
		orm.RegisterModel(new(models.User))
		orm.RegisterModel(new(models.MenuModel))
		orm.RegisterModel(new(captcha.MysqlCaptchaModel))
		// orm.RegisterModel(new(models.User))
}
