package controllers

import (
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/services"
		"github.com/WebGameLinux/cms/models"
		"github.com/astaxie/beego"
)

type ApiUserController interface {
		GetUserById()
}

type UserController struct {
		beego.Controller
}

func (this *UserController) URLMapping() {
		this.Mapping("Register", this.Register)
		this.Mapping("GetUserById", this.GetUserById)
}

// @router /v1/user/:id [get]
func (this *UserController) GetUserById() {
		var res services.ResultStruct
		if n, err := this.GetInt64(":id"); err == nil {
				res = services.GetUserBaseService().GetById(n)
		} else {
				res = services.NewParamsErrorResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
		}
		this.Data["json"] = res.Mapper()
		this.ServeJSON()
}

// @router /v1/user/register [post]
func (this *UserController) Register() {
		var (
				user = new(models.User)
				res  = services.NewSuccessResult(nil)
		)
		if err := this.ParseForm(user); err != nil {
				res.Set("code", enums.InvalidParams.Int())
				res.Set("message", enums.InvalidParams.WrapMsg(err))
		} else {
				res = services.GetUserRegisterService().RegisterByUser(user)
		}
		this.Data["json"] = res.Mapper()
		this.ServeJSON()
}
