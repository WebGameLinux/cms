package models

import (
		"errors"
		"github.com/WebGameLinux/cms/models/conditions"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/reflects"
		string2 "github.com/WebGameLinux/cms/utils/string"
		"github.com/WebGameLinux/cms/utils/types"
		"github.com/astaxie/beego/orm"
		"time"
)

// 用户操作封装
type UserWrapper struct {
		BaseWrapper
}

// 用户信息分页器
type UserPaginator struct {
		types.BasePaginator
}

func CreateUserPaginator() *UserPaginator {
		var paginator = new(UserPaginator)
		paginator.Provider = types.GetPaginatorProvider()
		paginator.Store("class", reflects.ClassName(paginator))
		return paginator
}

type UserModel interface {
		GetById(int64) *User
}

// 用户表名
func (user *User) TableName() string {
		return utils.GetTable("users")
}

// 默认用户封装实体构造
func NewUserWrapper() *UserWrapper {
		var wrapper = new(UserWrapper)
		wrapper.Model = new(User)
		return wrapper
}

// 获取用户封装实体
func GetUser(options ...interface{}) *UserWrapper {
		var wrapper = NewUserWrapper()
		WrapperInitOptions(wrapper, options...)
		return wrapper
}

// 用户名密码加密
func (wrapper *UserWrapper) Password(text string, options ...interface{}) string {
		if text == "" {
				return ""
		}
		if len(options) == 0 {
				return types.StrToPassword(text).Hash()
		}
		if v, ok := options[0].(*string2.PasswordHashOptions); ok {
				return types.StrToPassword(text).Hash(v)
		}
		if v, ok := options[0].(string2.PasswordHashOptions); ok {
				return types.StrToPassword(text).Hash(&v)
		}
		return types.StrToPassword(text).Hash()
}

// 验证密码是否相等处理
func (wrapper *UserWrapper) PasswordVerify(pass string, text string) bool {
		return types.StrToPassword(pass).Verify(text)
}

// 密码加密
func (wrapper *UserWrapper) PasswordHashed(text string) bool {
		if text == "" {
				return false
		}
		return types.StrToPassword(text).IsHashed()
}

// 验证用户登陆密码是否匹配
func (wrapper *UserWrapper) VerifyPasswordById(id int64, pass string) bool {
		var user = wrapper.GetById(id)
		if user == nil {
				return false
		}
		return wrapper.PasswordVerify(user.PasswordHash, pass)
}

//  通过ID 获取
func (wrapper *UserWrapper) GetById(id int64) *User {
		if id == 0 {
				return nil
		}
		user := new(User)
		user.Id = id
		err := wrapper.GetOrm().Read(user)
		if err == orm.ErrNoRows {
				return nil
		}
		return user
}

// 通过某一个字段获取用户信息
func (wrapper *UserWrapper) GetByKey(key string, value interface{}) *User {
		var user = new(User)
		if !wrapper.HasField(key) {
				return user
		}
		query, err := wrapper.GetQuery()
		if err != nil {
				wrapper.Error = err
				return user
		}
		query.Select(key).From(wrapper.Table()).Where(key + " = ?").
				OrderBy("id").SetModel(user).First(value)
		wrapper.Error = query.GetError()
		if wrapper.Error != nil {
				return nil
		}
		return user
}

// 用户信息更新
func (wrapper *UserWrapper) Update(id int64, data mapper.Mapper) *User {
		var user = wrapper.GetById(id)
		if user == nil {
				wrapper.Error = errors.New("用户不存在")
				return nil
		}
		var n int64
		// 获取对象属性集合
		m := reflects.GetItemsAllValues(user)
		if m == nil {
				wrapper.Error = errors.New("无权限更新")
				return nil
		}
		userMap := mapper.Mapper(m)
		// 属性集合 对比
		updateMap := userMap.Diff(data)
		if updateMap == nil || userMap.Len() == 0 {
				wrapper.Error = errors.New("无更新字段")
				return nil
		}
		// 更新对象属性
		if !mapper.SetByMap(user, updateMap) {
				wrapper.Error = errors.New("参数不符,无更新字段")
				return nil
		}
		// 更新数据库
		user.Version++
		user.UpdatedAt = time.Now()
		n, wrapper.Error = wrapper.GetOrm().Update(user, userMap.Keys()...)
		if wrapper.Error != nil && n > 0 {
				return nil
		}
		return user
}

// 创建
func (wrapper *UserWrapper) Create(user *User) int64 {
		if user == nil || user.PasswordHash == "" {
				return 0
		}
		if !wrapper.PasswordHashed(user.PasswordHash) {
				user.PasswordHash = wrapper.Password(user.PasswordHash)
		}
		id, err := wrapper.GetOrm().Insert(user)
		if err != nil {
				utils.Onerror(err)
				return 0
		}
		return id
}

// 获取一个
func (wrapper *UserWrapper) GetOne(user *User, columns ...string) *User {
		if user == nil {
				return user
		}
		wrapper.Error = wrapper.GetOrm().Read(user, columns...)
		return user
}

// 罗列
func (wrapper *UserWrapper) Lists(cond map[string]interface{}) []User {
		var users []User
		query, err := wrapper.GetQuery()
		if err != nil {
				wrapper.Error = err
				return users
		}
		pageOpt := conditions.MapCreatePageOptions(cond)
		query.Select("*").From(wrapper.Table()).
				SetModel(&users).
				OrderBy("id").
				Paginator(pageOpt.Page, pageOpt.Count, nil)
		wrapper.Error = query.GetError()
		return users
}

// 搜索
func (wrapper *UserWrapper) Search(cond map[string]interface{}) []User {
		var users []User
		query, err := wrapper.GetQuery()
		if err != nil {
				wrapper.Error = err
				return users
		}
		query.Select("*").From(wrapper.Table()).SetModel(&users).Get()
		wrapper.Error = query.GetError()
		return users
}
