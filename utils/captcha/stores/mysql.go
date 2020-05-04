package stores

import (
		"errors"
		"github.com/WebGameLinux/cms/utils/mapper"
		"log"
		"time"
)

type CaptchaMysqlStore struct {
		provider  MysqlStoreProvider
		Container map[string]interface{}
}

// mysql 服务提供器
type MysqlStoreProvider interface {
		Get(*CaptchaCodeDto) bool
		Set(*CaptchaCodeDto) bool
		Delete(*CaptchaCodeDto)
}

const (
		Code                  = "code"
		Id                    = "captcha_id"
		Provider              = "provider"
		Key                   = "key"
		Expire                = "expire"
		DefaultExpireDuration = 5 * time.Minute
)

// 验证码数据结构
type CaptchaCodeDto struct {
		Id        int64     `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Code      string    `orm:"column(code);size(30);index;description(验证码)" json:"code"`
		CaptchaId string    `orm:"column(captcha_id);size(128);description(验证ID);unique;" json:"captcha_id"`
		Key       string    `orm:"column(key);size(128);index;description(手机号,用户ID,其他键);null" json:"key"`
		ExpireAt  time.Time `orm:"column(expired_at);type(datetime);null;description(过期时间)" json:"expired_at"`
}

func NewCaptchaCodeDto() *CaptchaCodeDto {
		var dto = new(CaptchaCodeDto)
		return dto
}

func NewCaptchaMysqlStore(config ...mapper.Mapper) *CaptchaMysqlStore {
		if len(config) == 0 {
				config = append(config, make(mapper.Mapper))
		}
		var (
				store     = new(CaptchaMysqlStore)
				container = config[0]
		)
		store.Container = container
		provider := container.GetAny(Provider)
		if provider == nil {
				panic("miss mysql provider set")
		}
		if p, ok := provider.(MysqlStoreProvider); ok {
				store.provider = p
		}
		return store
}

func (this *CaptchaMysqlStore) SetMysqlProvider(provider MysqlStoreProvider) {
		if provider == nil {
				panic(errors.New("no provider set"))
		}
		this.provider = provider
}

func (this *CaptchaMysqlStore) Set(id string, value string) {
		var dto = new(CaptchaCodeDto)
		dto.Code = value
		dto.CaptchaId = id
		this.LoadConfig(dto)
		if !this.provider.Set(dto) {
				log.Fatal("create captcha code failed")
		}
		delete(this.Container, Key)
}

func (this *CaptchaMysqlStore) LoadConfig(dto *CaptchaCodeDto, keys ...[]string) {
		var (
				now = time.Now()
				m   = mapper.Mapper(this.Container)
		)
		if len(keys) == 0 {
				keys = append(keys, []string{Key, Expire})
		}
		if dto.Key == "" && this.include(Key, keys[0]) {
				dto.Key = m.Get(Key, "")
		}
		if dto.Key == "" {
				value := m.GetAny(Key)
				if value != nil {
						if handler, ok := value.(func(string, string) string); ok {
								dto.Key = handler(dto.CaptchaId, dto.Code)
						}
				}
		}
		t := time.Time{}
		if this.include(Expire, keys[0]) && (now.Before(dto.ExpireAt) || now.Equal(dto.ExpireAt) || t.Equal(dto.ExpireAt)) {
				d := m.GetDuration(Expire, time.Duration(0))
				if d <= time.Duration(0) {
						d = DefaultExpireDuration
				}
				dto.ExpireAt = now.Add(d)
		}
}

func (this *CaptchaMysqlStore) include(key string, keys []string) bool {
		for _, k := range keys {
				if k == key {
						return true
				}
		}
		return false
}

func (this *CaptchaMysqlStore) Get(id string, clear bool) string {
		var dto = new(CaptchaCodeDto)
		dto.CaptchaId = id
		this.LoadConfig(dto, []string{Key})
		if !this.provider.Get(dto) {
				return ""
		}
		if clear {
				defer this.provider.Delete(dto)
		}
		return dto.Code
}

func (this *CaptchaMysqlStore) Verify(id, answer string, clear bool) bool {
		return this.Get(id, clear) == answer
}
