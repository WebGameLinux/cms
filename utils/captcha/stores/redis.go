package stores

import (
		"errors"
		"github.com/WebGameLinux/cms/utils/mapper"
		"log"
		"time"
)

// redis 存储器
type CaptchaRedisStore struct {
		provider  RedisStoreProvider
		Container map[string]interface{}
}

// mysql 服务提供器
type RedisStoreProvider interface {
		Get(key string) string
		Set(key string, value string, expire ...int) bool
		Delete(key string)
}

// 数据结构对象
type RedisStoreDto struct {
		Code      string        `json:"code"`
		CaptchaId string        `json:"captcha_id"`
		Expire    time.Duration `json:"expire"`
}

func NewRedisCaptchaCodeDto() *RedisStoreDto {
		var dto = new(RedisStoreDto)
		return dto
}

func NewCaptchaRedisStore(config ...mapper.Mapper) *CaptchaRedisStore {
		if len(config) == 0 {
				config = append(config, make(mapper.Mapper))
		}
		var (
				store     = new(CaptchaRedisStore)
				container = config[0]
		)
		store.Container = container
		provider := container.GetAny(Provider)
		if provider == nil {
				panic("miss redis provider set")
		}
		if p, ok := provider.(RedisStoreProvider); ok {
				store.provider = p
		}
		return store
}

func (this *CaptchaRedisStore) SetMysqlProvider(provider RedisStoreProvider) {
		if provider == nil {
				panic(errors.New("no provider set"))
		}
		this.provider = provider
}

func (this *CaptchaRedisStore) Set(id string, value string) {
		var dto = new(RedisStoreDto)
		dto.Code = value
		dto.CaptchaId = id
		this.LoadConfig(dto)
		if !this.provider.Set(dto.CaptchaId, dto.Code, int(dto.Expire.Seconds())) {
				log.Fatal("create captcha code failed")
		}
}

func (this *CaptchaRedisStore) LoadConfig(dto *RedisStoreDto, keys ...[]string) {
		var m = mapper.Mapper(this.Container)
		if len(keys) == 0 {
				keys = append(keys, []string{Expire})
		}
		if this.include(Expire, keys[0]) {
				d := m.GetDuration(Expire, time.Duration(0))
				if d <= time.Duration(0) {
						d = DefaultExpireDuration
				}
				dto.Expire = d
		}
}

func (this *CaptchaRedisStore) include(key string, keys []string) bool {
		for _, k := range keys {
				if k == key {
						return true
				}
		}
		return false
}

func (this *CaptchaRedisStore) Get(id string, clear bool) string {
		code := this.provider.Get(id)
		if code == "" {
				return ""
		}
		if clear {
				defer this.provider.Delete(id)
		}
		return code
}

func (this *CaptchaRedisStore) Verify(id, answer string, clear bool) bool {
		return this.Get(id, clear) == answer
}
