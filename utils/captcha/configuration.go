package captcha

import (
		"github.com/mojocn/base64Captcha"
		"regexp"
)

const (
		ConfStoreKey         = "store"
		ConfTypeKey          = "type"
		ConfIdKey            = "id"
		ConfDriverKey        = "driver"
		DriverInstancePrefix = "driver:"
		StoreInstancePrefix  = "store:"
		DriverDigit          = "digit"
		DriverAudio          = "audio"
		DriverChinese        = "chinese"
		DriverEnglish        = "en"
		DriverString         = "string"
		DriverMath           = "math"
)


// 配置
type ConfigurationDto struct {
		Id          string `json:"id"`
		CaptchaType string `json:"captcha_type"`
		Driver      string `json:"driver"`
		Store       string `json:"store"`
}

// 驱动
type Drivers struct {
		DriverAudio   *base64Captcha.DriverAudio
		DriverString  *base64Captcha.DriverString
		DriverChinese *base64Captcha.DriverChinese
		DriverMath    *base64Captcha.DriverMath
		DriverDigit   *base64Captcha.DriverDigit
}

// 配置封装器
type ConfigurationWrapper struct {
		ConfigurationDto
		StoreInstance  base64Captcha.Store
		DriverInstance base64Captcha.Driver
}

type ConfigurationFace interface {
		GetId() string
		GetType() string
		GetDriver() string
		GetDriverInstance() base64Captcha.Driver
		GetStoreInstance() base64Captcha.Store
}

func NewConfiguration(data map[string]interface{}) ConfigurationFace {
		var m = new(ConfigurationWrapper)
		load(m, data)
		return m
}

func load(conf *ConfigurationWrapper, data map[string]interface{}) *ConfigurationWrapper {
		var (
				driverRegexp = regexp.MustCompile("^" + DriverInstancePrefix)
				storeRegexp  = regexp.MustCompile("^" + StoreInstancePrefix)
		)
		for key, v := range data {
				switch key {
				case ConfDriverKey:
						if str, ok := v.(string); ok {
								conf.Driver = str
								continue
						}
						if str, ok := v.(*string); ok {
								conf.Driver = *str
								continue
						}
						continue
				case ConfStoreKey:
						if str, ok := v.(string); ok {
								conf.Store = str
								continue
						}
						if str, ok := v.(*string); ok {
								conf.Driver = *str
								continue
						}
						continue
				case ConfTypeKey:
						if str, ok := v.(string); ok {
								conf.CaptchaType = str
								continue
						}
						if str, ok := v.(*string); ok {
								conf.CaptchaType = *str
								continue
						}
						continue
				case ConfIdKey:
						if str, ok := v.(string); ok {
								conf.Id = str
								continue
						}
						if str, ok := v.(*string); ok {
								conf.Id = *str
								continue
						}
						continue
				}
				if driverRegexp.MatchString(key) {
						if d, ok := v.(base64Captcha.Driver); ok {
								conf.DriverInstance = d
						}
						continue
				}
				if storeRegexp.MatchString(key) {
						if d, ok := v.(base64Captcha.Store); ok {
								conf.StoreInstance = d
						}
						continue
				}
		}
		return conf
}
