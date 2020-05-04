package controllers

import (
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/services"
		"github.com/astaxie/beego"
)

type ApiUserController interface {
		GetUserById()
}

type UserController struct {
		beego.Controller
}

func (this *UserController) URLMapping() {
		this.Mapping("GetUserById", this.GetUserById)
		this.Mapping("Register", this.Register)
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
		var res services.ResultStruct
		if n, err := this.GetInt64(":id"); err == nil {
				res = services.GetUserBaseService().GetById(n)
		} else {
				res = services.NewParamsErrorResult(3000, err.Error())
		}
		this.Data["json"] = res.Mapper()
		this.ServeJSON()
}

