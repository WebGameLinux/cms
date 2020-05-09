package services

import (
		"encoding/json"
		"errors"
		"fmt"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/utils/mapper"
		string2 "github.com/WebGameLinux/cms/utils/string"
		"github.com/astaxie/beego/cache"
		"time"
)

type TokenService interface {
		Verify(token string) bool
		Store(data map[string]interface{}, expire ...time.Duration) string
		Get(token string) (mapper.Mapper, bool)
		Load(token string, v interface{}) bool
		SetStore(store string) TokenService
		GetError() error
}

type UserAuthTokenService struct {
		store   string
		expire  time.Duration
		storage cache.Cache
		err     error
}


var tokenService TokenService

func GetTokenService() TokenService {
		if tokenService == nil {
				var service = new(UserAuthTokenService)
				service.store = enums.TokenStore
				service.expire = 6 * time.Minute
				tokenService = service
		}
		return tokenService
}

func (this *UserAuthTokenService) Dispatch(token string,model...interface{}) {

}

func (this *UserAuthTokenService) keep(token string)  {
		this.GetStorage()
}

func (this *UserAuthTokenService) SetStore(store string) TokenService {
		if this.store != store {
				this.storage = nil
		}
		this.store = store
		return this
}

func (this *UserAuthTokenService) Verify(token string) bool {
		return this.GetStorage().IsExist(token)
}

func (this *UserAuthTokenService) Store(data map[string]interface{}, expire ...time.Duration) string {
		if len(expire) == 0 {
				expire = append(expire, this.expire)
		}
		var token string
		storeStr := this.value(data)
		if storeStr == "" {
				return ""
		}
		token = this.token(storeStr)
		this.err = this.GetStorage().Put(token, storeStr, expire[0])
		return token
}

func (this *UserAuthTokenService) value(data mapper.Mapper) string {
		var (
				buf []byte
				err error
		)
		data["publish_at"] = time.Now().String()
		if buf, err = json.Marshal(data); err != nil {
				this.err = err
				return ""
		}
		return string(buf)
}

func (this *UserAuthTokenService) GetError() error {
		var err = this.err
		if this.err != nil {
				this.err = nil
		}
		return err
}

func (this *UserAuthTokenService) token(data string) string {
		return string2.Md5(data)
}

func (this *UserAuthTokenService) Get(token string) (mapper.Mapper, bool) {
		val := this.GetStorage().Get(token)
		if val == nil || val == "" {
				return nil, false
		}
		if buf, ok := val.([]byte); ok {
				val = string(buf)
		}
		if data, ok := val.(string); ok {
				m := this.mapper(data)
				if m != nil && len(m) > 0 {
						return m, true
				}
		}
		this.err = errors.New("unknown token data type " + fmt.Sprintf("%+v", val))
		return nil, false
}

func (this *UserAuthTokenService) Load(token string, v interface{}) bool {
		if v == nil {
				return false
		}
		val := this.GetStorage().Get(token)
		if val == nil || val == "" {
				return false
		}
		if buf, ok := val.([]byte); ok {
				return this.bind(string(buf), v)
		}
		if data, ok := val.(string); ok {
				return this.bind(data, v)
		}
		this.err = errors.New("unknown token data type " + fmt.Sprintf("%+v", val))
		return false
}

func (this *UserAuthTokenService) bind(data string, v interface{}) bool {
		if data == "" {
				return false
		}
		this.err = json.Unmarshal([]byte(data), v)
		return true
}

func (this *UserAuthTokenService) mapper(data string) mapper.Mapper {
		if data == "" {
				return nil
		}
		var m = make(mapper.Mapper)
		this.err = json.Unmarshal([]byte(data), &m)
		return m
}

func (this *UserAuthTokenService) GetStorage() cache.Cache {
		if this.storage == nil {
				this.storage = GetCacheManagerService().Store(this.store)
		}
		return this.storage
}
