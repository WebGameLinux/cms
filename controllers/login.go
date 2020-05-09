package controllers

import (
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/repositories"
		"github.com/WebGameLinux/cms/dto/request"
		"github.com/WebGameLinux/cms/dto/response"
		"github.com/astaxie/beego"
		"strconv"
		"time"
)

type LoginController struct {
		BaseController
}

type ApiLoginController interface {
		Login()
}

var loginController beego.ControllerInterface

func GetLoginController() beego.ControllerInterface {
		if loginController == nil {
				loginController = new(LoginController)
		}
		return loginController
}

func (this *LoginController) URLMapping() {
		this.Mapping("Login", this.Login)
}

// @router /v1/login [post]
func (this *LoginController) Login() {
		var res common.JsonResponseInterface
		params := this.GetParamsMapper()
		if params == nil {
				res = common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg())
		} else {
				dto := new(request.LoginParamsDto)
				res = repositories.GetUserRepository().Login(dto.Load(params))
		}
		this.cookie(res)
		this.ApiResponse(res)
}

func (this *LoginController) cookie(res common.JsonResponseInterface) {
		if !res.IsSuccess() {
				return
		}
		data := res.Item()
		if v, ok := data.(*response.LoginFilterRespJson); ok {
				expire := 7 * 24 * time.Hour
				this.Ctx.SetCookie(enums.UserId, strconv.Itoa(int(v.User.Id)), expire, "/")
				this.Ctx.SetCookie(enums.AuthToken, v.Auth, expire, "/")
				d, _ := v.MarshalJSON()
				this.Ctx.SetCookie("user", string(d), expire, "/")
		}
}
