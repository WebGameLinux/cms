package services

import (
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/repositories"
		"github.com/WebGameLinux/cms/utils/mapper"
)

type UserService struct {
		repo repositories.UserRepository
}

type UserBaseService interface {
		GetById(int64) ResultStruct
		VerifyPasswordById(id int64, pass string, encode ...bool) bool
		CreateUser(map[string]interface{}) ResultStruct
}

func NewUserService() *UserService {
		return &UserService{repo: repositories.GetUserBaseRepository()}
}

func GetUserBaseService(options ...interface{}) UserBaseService {
		return &UserService{repo: repositories.GetUserBaseRepository(options...)}
}

// 通过uid 获取用户信息
func (user *UserService) GetById(id int64) ResultStruct {
		var data = user.repo.GetById(id)
		if data == nil {
				return NewEmptyResult(enums.RecordNotExists.Int(), enums.RecordNotExists.WrapMsg(user.repo.GetError()))
		}
		return NewSuccessResult(mapper.NewKvMap("user", data))
}

// 创建用户
func (user *UserService) CreateUser(data map[string]interface{}) ResultStruct {
		id := user.repo.CreateByMap(data)
		err := user.repo.GetError()
		if id <= 0 || err != nil {
				return NewEmptyResult(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg(err))
		}
		return NewSuccessResult(user.GetById(id).Item(), enums.SUCCESS.String())
}

// 验证登陆密码
func (user *UserService) VerifyPasswordById(id int64, pass string, encode ...bool) bool {
		if len(encode) == 0 {
				encode = append(encode, false)
		}
		return user.repo.VerifyPasswordById(id, pass, encode[0])
}

// 更新用户记录
func (user *UserService) UpdateById(id int64, data Map) ResultStruct {
		res := user.repo.Update(id, data)
		err := user.repo.GetError()
		if res == nil || err != nil {
				return NewEmptyResult(enums.UpdateModelField.Int(), enums.UpdateModelField.WrapMsg(err))
		}
		return NewSuccessResult(res, enums.SUCCESS.String())
}
