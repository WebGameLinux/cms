package services

import (
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/validator"
)

type RegisterServiceInterface interface {
		RegisterByUser(user *models.User) common.JsonResponseInterface
		RegisterByMap(data map[string]interface{}) common.JsonResponseInterface
}

var userRegisterService RegisterServiceInterface

type UserRegisterService struct {
		model *models.UserWrapper
}

func NewRegisterService() *UserRegisterService {
		service := new(UserRegisterService)
		service.model = models.GetUser()
		return service
}

func GetUserRegisterService() RegisterServiceInterface {
		if userRegisterService == nil {
				userRegisterService = NewRegisterService()
		}
		return userRegisterService
}

func (this *UserRegisterService) CheckEmail(email string) bool {
		return this.model.GetByKey(enums.Email, email, false) != nil
}

func (this *UserRegisterService) CheckUserName(name string) bool {
		return this.model.GetByKey(enums.UserName, name, false) != nil
}

func (this *UserRegisterService) CheckMobile(mobile string) bool {
		return this.model.GetByKey(enums.Mobile, mobile, false) != nil
}

func (this *UserRegisterService) Check(user *models.User) common.JsonResponseInterface {
		if user == nil {
				return common.Warn(enums.ParamEmpty.Int(), enums.ParamEmpty.WrapMsg())
		}
		if this.CheckEmail(user.Email) {
				return common.Warn(enums.KeyExists.Int(), enums.KeyExists.Replace("邮箱已存在"))
		}
		if this.CheckMobile(user.Mobile) {
				return common.Warn(enums.KeyExists.Int(), enums.KeyExists.Replace("手机号已存在"))
		}
		if this.CheckUserName(user.UserName) {
				return common.Warn(enums.KeyExists.Int(), enums.KeyExists.Replace("用户名已存在"))
		}
		return common.Success(nil)
}

func (this *UserRegisterService) RegisterByMap(data map[string]interface{}) common.JsonResponseInterface {
		var user = new(models.User)
		m := mapper.Mapper(data)
		m.ReName(enums.Password, enums.PasswordHash)
		if err := user.LoadByMap(data); err != nil {
				return common.Warn(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
		}
		if _, err := user.Valid(); err != nil {
				info := validator.TranslateError(err)
				return common.Warn(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(info))
		}
		check := this.Check(user)
		if !check.IsSuccess() {
				return check
		}
		id := this.model.Create(user)
		if id <= 0 {
				return common.Warn(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg())
		}
		return GetLoginService().LoginByUid(id)
}

func (this *UserRegisterService) RegisterByUser(user *models.User) common.JsonResponseInterface {
		if _, err := user.Valid(); err != nil {
				info := validator.TranslateError(err)
				return common.Warn(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(info))
		}
		check := this.Check(user)
		if !check.IsSuccess() {
				return check
		}
		id := this.model.Create(user)
		if id <= 0 {
				return common.Warn(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg())
		}
		return GetLoginService().LoginByUid(id)
}
