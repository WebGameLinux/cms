package controllers

import (
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/repositories"
		"github.com/astaxie/beego"
)

type ApiUserController interface {
		GetUserById()
		Register()
}

type UserController struct {
		BaseController
}

var userController beego.ControllerInterface

func GetUserController() beego.ControllerInterface {
		if userController == nil {
				userController = new(UserController)
		}
		return userController
}

func (this *UserController) URLMapping() {
		this.Mapping("Register", this.Register)
		this.Mapping("GetUserById", this.GetUserById)
}

// @router /v1/user/:id [get]
func (this *UserController) GetUserById() {
		var res common.JsonResponseInterface
		if n, err := this.GetInt64(":id"); err == nil {
				res = repositories.GetUserRepository().GetById(n)
		} else {
				res = common.ParamError(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
		}
		this.Data["json"] = res.Mapper()
		this.ServeJSON()
}

// @router /v1/register [post]
func (this *UserController) Register() {
		var res = common.Success(nil)
		data := this.GetJsonMapper()
		if data == nil {
				res.Set("code", enums.InvalidParams.Int())
				res.Set("message", enums.InvalidParams.WrapMsg())
		} else {
				res = repositories.GetUserRepository().RegisterByMap(data)
		}
		this.ApiResponse(res)
}
