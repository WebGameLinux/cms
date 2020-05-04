package services

import (
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/repositories"
		"github.com/WebGameLinux/cms/models"
)

type RegisterService interface {
		RegisterByMap(data map[string]interface{}) ResultStruct
		RegisterByUser(user *models.User) ResultStruct
}

var userRegisterService RegisterService

type UserRegisterService struct {
		repo repositories.UserRepository
}

func GetUserRegisterService() RegisterService {
		if userRegisterService == nil {
				service := new(UserRegisterService)
				service.repo = repositories.GetUserRepository()
				userRegisterService = service
		}
		return userRegisterService
}

func (this *UserRegisterService) RegisterByMap(data map[string]interface{}) ResultStruct {
		var user = new(models.User)
		if err := user.LoadByMap(data); err != nil {
				return NewEmptyResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
		}
		if verify, err := user.Valid(); err != nil {
				if verify == nil {
						return NewEmptyResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
				}
				return NewEmptyResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(verify.Errors[0]))
		}
		id := this.repo.Create(user)
		if id <= 0 {
				return NewEmptyResult(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg())
		}
		return NewSuccessResult(this.repo.GetById(id))
}

func (this *UserRegisterService) RegisterByUser(user *models.User) ResultStruct {
		if verify, err := user.Valid(); err != nil {
				if verify == nil {
						return NewEmptyResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(err))
				}
				return NewEmptyResult(enums.InvalidParams.Int(), enums.InvalidParams.WrapMsg(verify.Errors[0]))
		}
		id := this.repo.Create(user)
		if id <= 0 {
				return NewEmptyResult(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg())
		}
		return NewSuccessResult(this.repo.GetById(id))
}
