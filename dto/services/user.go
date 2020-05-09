package services

import (
		"errors"
		enums2 "github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/mapper"
)

var (
		userService *UserService
)

// 获取用户
type UserServiceInterface interface {
		GetById(int64) *models.User
		Password(int64) string
		GetError() error
		VerifyPasswordById(id int64, pass string, encode bool) bool
		Exists(string, interface{}, ...bool) bool
		GetUserProperties(id int64, keys []string) mapper.Mapper
		Update(id int64, data map[string]interface{}) *models.User
		CreateByMap(data map[string]interface{}) int64
		Create(data *models.User) int64
		VerifyPassword(pass string, txt string) bool
		GetByUserName(name string) *models.User
		GetByEmail(email string) *models.User
		GetByMobile(mobile string) *models.User
}

type UserService struct {
		model *models.UserWrapper
}

func NewUserService() *UserService {
		service := new(UserService)
		service.model = models.GetUser()
		return service
}

func GetUserService() UserServiceInterface {
		if nil == userService {
				userService = NewUserService()
		}
		return userService
}

func (this *UserService) Exists(key string, v interface{}, softDelete ...bool) bool {
		if this.model.GetByKey(key, v, softDelete...) != nil {
				return true
		}
		return false
}

func (this *UserService) GetUserProperties(id int64, keys []string) mapper.Mapper {
		panic("implement me")
}

func (this *UserService) VerifyPasswordById(id int64, pass string, encode bool) bool {
		if pass == "" {
				return false
		}
		var data = this.model.GetById(id)
		if data == nil {
				return false
		}
		if data.PasswordHash == "" {
				return false
		}
		return true
}

func (this *UserService) Lists() {

}

func (this *UserService) GetOne(user *models.User) bool {
		if err := this.model.GetOrm().Read(user); err == nil {
				return true
		}
		return false
}

func (this *UserService) GetById(id int64) *models.User {
		return this.model.GetById(id)
}

func (this *UserService) Update(id int64, data map[string]interface{}) *models.User {
		// 字段重命名
		if v, ok := data["password"]; ok {
				m := mapper.Mapper(data)
				m.ReName("password", "passwordHash")
				if v == nil || v == "" {
						m.Delete("password")
				}
		}
		return this.model.Update(id, data)
}

func (this *UserService) Password(id int64) string {
		model := this.model.GetById(id)
		if model == nil {
				return ""
		}
		return model.PasswordHash
}

func (this *UserService) GetError() error {
		return this.model.GetError()
}

func (this *UserService) CreateByMap(data map[string]interface{}) int64 {
		model := new(models.User)
		m := mapper.Mapper(data)
		m.ReName(enums2.Password, enums2.PasswordHash)
		if mapper.SetByMap(model, m) {
				this.model.Error = errors.New("数据格式不匹配")
				return 0
		}
		id := this.model.Create(model)
		return id
}

// 通过user 模型创建
func (this *UserService) Create(data *models.User) int64 {
		id := this.model.Create(data)
		return id
}

func (this *UserService) GetByMobile(mobile string) *models.User {
		return this.model.GetByKey(enums2.Mobile, mobile)
}

func (this *UserService) GetByEmail(email string) *models.User {
		return this.model.GetByKey(enums2.Email, email)
}

func (this *UserService) GetByUserName(name string) *models.User {
		return this.model.GetByKey(enums2.UserName, name)
}

func (this *UserService) VerifyPassword(pass string, txt string) bool {
		return this.model.PasswordVerify(pass, txt)
}
