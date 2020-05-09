package services

import (
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/request"
		"github.com/WebGameLinux/cms/dto/response"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/validator"
		"regexp"
)

type LoginServiceInterface interface {
		Login(params *request.LoginParamsDto) common.JsonResponseInterface
		LoginByUser(user *models.User) common.JsonResponseInterface
		LoginByUid(uid int64) common.JsonResponseInterface
		LoginByMobile(mobile string) common.JsonResponseInterface
		LoginByEmail(email string) common.JsonResponseInterface
}

type UserLoginService struct {
		model *models.UserWrapper
}

var loginService LoginServiceInterface

func NewLoginService() *UserLoginService {
		service := new(UserLoginService)
		service.model = models.GetUser()
		return service
}

func GetLoginService() LoginServiceInterface {
		if loginService == nil {
				loginService = NewLoginService()
		}
		return loginService
}

func (this *UserLoginService) Login(params *request.LoginParamsDto) common.JsonResponseInterface {
		if params == nil {
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg("登陆参数异常"))
		}
		if _, err := params.Valid(); err != nil {
				info := validator.TranslateError(err)
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg(info))
		}
		typ := this.CheckAccountType(params.Account)
		if typ == "" {
				return common.Warn(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg("未知登陆方式"))
		}
		var user *models.User
		switch typ {
		case enums.UserName:
				user = this.model.GetByKey(enums.UserName, params.Account)
		case enums.Email:
				user = this.model.GetByKey(enums.Email, params.Account)
		case enums.Mobile:
				user = this.model.GetByKey(enums.Mobile, params.Account)
		}
		if user == nil {
				return common.Warn(enums.ErrorUserNotExist.Int(), enums.ErrorUserNotExist.WrapMsg("请检查账号"))
		}
		if !this.model.PasswordVerify(user.PasswordHash, params.Password) {
				return common.Warn(enums.ErrorUserLogin.Int(), enums.ErrorUserLogin.WrapMsg())
		}
		token := GetTokenService().Store(user.Filter())
		if token == "" {
				return common.Warn(enums.CreateTokenFailed.Int(), enums.CreateTokenFailed.WrapMsg(GetTokenService().GetError()))
		}

		return this.Success(user, token)
}

func (this *UserLoginService) LoginByUid(uid int64) common.JsonResponseInterface {
		user := this.model.GetById(uid)
		if user == nil {
				return common.Warn(enums.ErrorUserNotExist.Int(), enums.ErrorUserNotExist.WrapMsg())
		}
		token := GetTokenService().Store(user.Filter())
		if token == "" {
				return common.Warn(enums.CreateTokenFailed.Int(), enums.CreateTokenFailed.WrapMsg(GetTokenService().GetError()))
		}
		return this.Success(user, token)
}

func (this *UserLoginService) LoginByUser(user *models.User) common.JsonResponseInterface {
		if user == nil {
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg())
		}
		if user.Id == 0 || user.Mobile == "" || user.PasswordHash == "" || user.UserName == "" {
				return common.Warn(enums.ErrorUserNotExist.Int(), enums.ErrorUserNotExist.WrapMsg())
		}
		token := GetTokenService().Store(user.Filter())
		if token == "" {
				return common.Warn(enums.CreateTokenFailed.Int(), enums.CreateTokenFailed.WrapMsg(GetTokenService().GetError()))
		}
		return this.Success(user, token)
}

func (this *UserLoginService) LoginByMobile(mobile string) common.JsonResponseInterface {
		if mobile == "" {
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg())
		}
		user := this.model.GetByKey(enums.Mobile, mobile)
		if user == nil {
				return common.Warn(enums.ErrorUserNotExist.Int(), enums.ErrorUserNotExist.WrapMsg())
		}
		token := GetTokenService().Store(user.Filter())
		if token == "" {
				return common.Warn(enums.CreateTokenFailed.Int(), enums.CreateTokenFailed.WrapMsg(GetTokenService().GetError()))
		}
		return this.Success(user, token)
}

func (this *UserLoginService) LoginByEmail(email string) common.JsonResponseInterface {
		if email == "" {
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg())
		}
		user := this.model.GetByKey(enums.Email, email)
		if user == nil {
				return common.Warn(enums.ErrorUserNotExist.Int(), enums.ErrorUserNotExist.WrapMsg())
		}
		token := GetTokenService().Store(user.Filter())
		if token == "" {
				return common.Warn(enums.CreateTokenFailed.Int(), enums.CreateTokenFailed.WrapMsg(GetTokenService().GetError()))
		}
		return this.Success(user, token)
}

func (this *UserLoginService) CheckAccountType(account string) string {
		var (
				mobile = regexp.MustCompile(`1[2-9]{10}`)
				email  = regexp.MustCompile(`\w{1,100}@\w{1,50}.\w{1,5}`)
		)
		if mobile.MatchString(account) {
				return enums.Mobile
		}
		if email.MatchString(account) {
				return enums.Email
		}
		return enums.UserName
}

func (this *UserLoginService) Success(user *models.User, token string) common.JsonResponseInterface {
		data := new(response.LoginRespJson)
		data.User = user
		data.Auth = token
		return common.Success(data, "登陆成功")
}
