package captcha

import (
		"github.com/WebGameLinux/cms/utils/captcha/stores"
)

type RedisCaptchaProvider struct {
		
}

type RedisCaptchaWrapper interface {
		stores.RedisStoreProvider
}

type RedisCaptchaModel struct {
		stores.RedisStoreDto
}

func NewRedisCaptchaProvider() *RedisCaptchaProvider {
    var provider = new(RedisCaptchaProvider)
    return provider
}

func (this * RedisCaptchaProvider) Get(key string) string {
		panic("implement me")
}

func (this * RedisCaptchaProvider) Set(key string, value string, expire ...int) bool {
		panic("implement me")
}

func (this * RedisCaptchaProvider) Delete(key string) {
		panic("implement me")
}