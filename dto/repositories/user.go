package repositories

import (
		"errors"
		"github.com/WebGameLinux/cms/dto/transforms"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/models/enums"
		"github.com/WebGameLinux/cms/utils/mapper"
		"sync"
)

var (
		once                       sync.Once
		instanceUserBaseRepository *UserBaseRepository
)

// 获取用户
type UserRepository interface {
		GetById(int64) *models.User
		Password(int64) string
		GetError() error
		VerifyPasswordById(id int64, pass string, encode bool) bool
		Exists(string, interface{}) bool
		GetUserProperties(id int64, keys []string) mapper.Mapper
		Update(id int64, data map[string]interface{}) *models.User
		CreateByMap(data map[string]interface{}) int64
		Create(data *models.User) int64
}

type UserBaseRepository struct {
		Model *models.UserWrapper
}

func GetUserRepository() UserRepository {
		if nil == instanceUserBaseRepository {
				instanceUserBaseRepository = GetUserBaseRepository()
		}
		return instanceUserBaseRepository
}

func GetUserBaseRepository(options ...interface{}) *UserBaseRepository {
		var (
				v    = ""
				ok   bool
				repo = new(UserBaseRepository)
		)
		if len(options) == 0 {
				repo.Model = models.GetUser()
				return repo
		}
		name := options[0]
		if v, ok = name.(string); ok {
				v = ""
		}
		repo.Model = models.GetUser(v)
		return repo
}

func (user *UserBaseRepository) Exists(string, interface{}) bool {
		panic("implement me")
}

func (user *UserBaseRepository) GetUserProperties(id int64, keys []string) mapper.Mapper {
		panic("implement me")
}

func (user *UserBaseRepository) VerifyPasswordById(id int64, pass string, encode bool) bool {
		if pass == "" {
				return false
		}
		var data = user.Model.GetById(id)
		if data == nil {
				return false
		}
		if data.PasswordHash == "" {
				return false
		}
		return true
}

func (user *UserBaseRepository) Lists() {

}

func (user *UserBaseRepository) GetOne(conditions map[string]interface{}) *models.User {
		return nil
}

func (user *UserBaseRepository) GetById(id int64) *models.User {
		return user.Model.GetById(id)
}

func (user *UserBaseRepository) Update(id int64, data map[string]interface{}) *models.User {
		// 字段重命名
		if v, ok := data["password"]; ok {
				data["password_hash"] = v
				delete(data, "password_hash")
		}
		// gender 类型转换
		if v, ok := data["gender"]; ok {
				if n, ok := v.(int); ok {
						data["gender"] = enums.ParseInt(n)
				}
		}
		return user.Model.Update(id, data)
}

func (user *UserBaseRepository) Password(id int64) string {
		model := user.Model.GetById(id)
		if model == nil {
				return ""
		}
		return model.PasswordHash
}

func (user *UserBaseRepository) GetError() error {
		return user.Model.GetError()
}

func (user *UserBaseRepository) CreateByMap(data map[string]interface{}) int64 {
		model := new(models.User)
		m := mapper.Mapper(data)
		m.ReName("password", "password_hash")
		m.Transform("gender", transforms.TransformMapGender)
		if mapper.SetByMap(model, m) {
				user.Model.Error = errors.New("数据格式不匹配")
				return 0
		}
		id := user.Model.Create(model)
		return id
}

// 通过user 模型创建
func (user *UserBaseRepository) Create(data *models.User) int64 {
		id := user.Model.Create(data)
		return id
}
