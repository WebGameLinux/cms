package repositories

import (
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/request"
		"github.com/WebGameLinux/cms/dto/response"
		"github.com/WebGameLinux/cms/dto/services"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/mapper"
)

type UserRepositoryInterface interface {
		GetById(int64) common.JsonResponseInterface
		VerifyPasswordById(id int64, pass string, encode ...bool) bool
		CreateUser(map[string]interface{}) common.JsonResponseInterface
		Login(params *request.LoginParamsDto) common.JsonResponseInterface
		LoginByUser(user *models.User) common.JsonResponseInterface
		LoginByUid(uid int64) common.JsonResponseInterface
		LoginByMobile(mobile string) common.JsonResponseInterface
		LoginByEmail(email string) common.JsonResponseInterface
		RegisterByUser(user *models.User) common.JsonResponseInterface
		RegisterByMap(data map[string]interface{}) common.JsonResponseInterface
}

type UserRepository struct {
		service services.UserServiceInterface
}

var userRepository *UserRepository

func NewUserRepository() *UserRepository {
		return &UserRepository{
				service: services.GetUserService(),
		}
}

func GetUserRepository() UserRepositoryInterface {
		if userRepository == nil {
				userRepository = NewUserRepository()
		}
		return userRepository
}

// 通过uid 获取用户信息
func (this *UserRepository) GetById(id int64) common.JsonResponseInterface {
		var data = this.service.GetById(id)
		if data == nil {
				return common.Warn(enums.RecordNotExists.Int(), enums.RecordNotExists.WrapMsg(this.service.GetError()))
		}
		return common.Success(mapper.NewKvMap("user", data.Filter()))
}

// 创建用户
func (this *UserRepository) CreateUser(data map[string]interface{}) common.JsonResponseInterface {
		id := this.service.CreateByMap(data)
		err := this.service.GetError()
		if id <= 0 || err != nil {
				return common.Warn(enums.CreateRecordField.Int(), enums.CreateRecordField.WrapMsg(err))
		}
		return common.Success(this.GetById(id).Item(), enums.SUCCESS.String())
}

// 验证登陆密码
func (this *UserRepository) VerifyPasswordById(id int64, pass string, encode ...bool) bool {
		if len(encode) == 0 {
				encode = append(encode, false)
		}
		return this.service.VerifyPasswordById(id, pass, encode[0])
}

// 更新用户记录
func (this *UserRepository) UpdateById(id int64, data map[string]interface{}) common.JsonResponseInterface {
		res := this.service.Update(id, data)
		err := this.service.GetError()
		if res == nil || err != nil {
				return common.Warn(enums.UpdateModelField.Int(), enums.UpdateModelField.WrapMsg(err))
		}
		return common.Success(res, enums.SUCCESS.String())
}

func (this *UserRepository) Login(params *request.LoginParamsDto) common.JsonResponseInterface {
		return this.LoginResponseFilter(services.GetLoginService().Login(params))
}

func (this *UserRepository) LoginByUser(user *models.User) common.JsonResponseInterface {
		return this.LoginResponseFilter(services.GetLoginService().LoginByUser(user))
}

func (this *UserRepository) LoginByUid(uid int64) common.JsonResponseInterface {
		return this.LoginResponseFilter(services.GetLoginService().LoginByUid(uid))
}

func (this *UserRepository) LoginByMobile(mobile string) common.JsonResponseInterface {
		return this.LoginResponseFilter(services.GetLoginService().LoginByMobile(mobile))
}

func (this *UserRepository) LoginByEmail(email string) common.JsonResponseInterface {
		return this.LoginResponseFilter(services.GetLoginService().LoginByEmail(email))
}

func (this *UserRepository) RegisterByUser(user *models.User) common.JsonResponseInterface {
		return services.GetUserRegisterService().RegisterByUser(user)
}

func (this *UserRepository) RegisterByMap(data map[string]interface{}) common.JsonResponseInterface {
		return services.GetUserRegisterService().RegisterByMap(data)
}

func (this *UserRepository) LoginResponseFilter(responseInterface common.JsonResponseInterface) common.JsonResponseInterface {
		if !responseInterface.IsSuccess() {
				return responseInterface
		}
		data := responseInterface.Item()
		if v, ok := data.(*response.LoginRespJson); ok {
				d := new(response.LoginFilterRespJson)
				responseInterface.Set("data", d.Init(v))
		}
		return responseInterface
}
