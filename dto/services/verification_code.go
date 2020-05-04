package services

import (
		captcha2 "github.com/WebGameLinux/cms/models/captcha"
		"github.com/WebGameLinux/cms/utils/captcha"
		"github.com/WebGameLinux/cms/utils/captcha/drivers"
		"github.com/WebGameLinux/cms/utils/captcha/stores"
		"github.com/mojocn/base64Captcha"
		"strings"
		"time"
)

const (
		TypeMobile = "mobile"
		TypeImage  = "image"
		OptLength  = "length"
		OptExpire  = "expire"
)

// 验证码服务
type VerificationService interface {
		Mobile(mobile string, options ...*MobileOption) *captcha.Wrapper
		Image(name ...string) *captcha.Wrapper
		Verify(id string, code string, ty string, clean ...bool) bool
		SetStore(name string, instance base64Captcha.Store)
		SetDriver(name string, instance base64Captcha.Driver)
}

// 验证码服务
type VerificationCodeService struct {
		Stores           map[string]base64Captcha.Store
		DriversContainer map[string]base64Captcha.Driver
}

type MobileOption struct {
		Length int           // 文字数量
		Expire time.Duration // 存活时长
}

func NewMobileOption(args ...interface{}) *MobileOption {
		var opt = new(MobileOption)
		if len(args) == 0 {
				opt.Length = 6
				opt.Expire = time.Minute * 5
				return opt
		}
		for _, v := range args {
				if length, ok := v.(int); ok && length > 0 {
						opt.Length = length
				}
				if expire, ok := v.(time.Duration); ok && expire > time.Duration(0) {
						opt.Expire = expire
				}
		}
		if opt.Expire <= time.Duration(0) {
				opt.Expire = time.Minute * 5
		}
		if opt.Length <= 0 {
				opt.Length = 6
		}
		return opt
}

func GetVerificationService() VerificationService {
		var service = NewVerificationService()
		loadVerificationService(service)
		return service
}

func NewVerificationService() *VerificationCodeService {
		var service = new(VerificationCodeService)
		service.Stores = make(map[string]base64Captcha.Store)
		service.DriversContainer = make(map[string]base64Captcha.Driver)
		return service
}

func loadVerificationService(service *VerificationCodeService) *VerificationCodeService {
		service.SetDriver(TypeMobile, drivers.NewCaptchaMobileDriver("", 6, ""))
		m := make(map[string]interface{})
		m[stores.Expire] = 5 * time.Minute
		model := captcha2.NewMysqlCaptchaWrapper()
		m[stores.Provider] = model
		service.SetStore(TypeMobile, stores.NewCaptchaMysqlStore(m))
		service.SetDriver(TypeImage, base64Captcha.DefaultDriverDigit)
		service.SetStore(TypeImage, base64Captcha.DefaultMemStore)
		return service
}

// 验证码生成封装
func (this *VerificationCodeService) Code(args ...interface{}) ResultStruct {
		return nil
}

func (this *VerificationCodeService) Mobile(mobile string, options ...*MobileOption) *captcha.Wrapper {
		if len(options) == 0 {
				options = append(options, NewMobileOption())
		}
		var (
				opt = options[0]
				m   = this.getMap(TypeMobile, mobile)
		)
		m[TypeMobile] = mobile
		m[OptLength] = opt.Length
		m[OptExpire] = opt.Expire
		return captcha.NewCaptcha(m)
}

func (this *VerificationCodeService) Image(name ...string) *captcha.Wrapper {
		if len(name) == 0 {
				name = append(name, TypeImage)
		}
		return captcha.NewCaptcha(this.getMap(name[0], ""))
}

func (this *VerificationCodeService) GetStore(name string) base64Captcha.Store {
		if v, ok := this.Stores[name]; ok {
				return v
		}
		return this.getDefaultStore()
}

func (this *VerificationCodeService) GetDriver(name string, args ...string) base64Captcha.Driver {
		if v, ok := this.DriversContainer[name]; ok {
				return v
		}
		if name == TypeMobile {
				if len(args) != 0 {
						return this.getMobileDriver(args[0], 6, 5*time.Minute)
				}
				return this.getMobileDriver("", 6, 5*time.Minute)
		}
		return this.getDefaultDriver()
}

func (this *VerificationCodeService) getDefaultDriver() base64Captcha.Driver {
		return base64Captcha.DefaultDriverDigit
}

func (this *VerificationCodeService) getMobileDriver(mobile string, length int, d time.Duration) base64Captcha.Driver {
		return drivers.NewCaptchaMobileDriver(mobile, length, "")
}

func (this *VerificationCodeService) getDefaultStore() base64Captcha.Store {
		return base64Captcha.DefaultMemStore
}

func (this *VerificationCodeService) getMap(name string, mobile string) map[string]interface{} {
		m := make(map[string]interface{})
		driver := this.GetDriver(name, mobile)
		if mobile != "" {
				if d, ok := driver.(*drivers.CaptchaMobileDriver); ok {
						d.Mobile = mobile
				}
		}
		m[this.DriverName(name)] = driver
		store := this.GetStore(name)
		// 注入mobile
		if p, ok := store.(*stores.CaptchaMysqlStore); ok && mobile != "" {
				p.Container[stores.Key] = mobile
		}
		m[this.StoreName(name)] = store
		return m
}

func (this *VerificationCodeService) StoreName(name string) string {
		if strings.Contains(name, captcha.StoreInstancePrefix) {
				return name
		}
		return captcha.StoreInstancePrefix + name
}

func (this *VerificationCodeService) DriverName(name string) string {
		if strings.Contains(name, captcha.DriverInstancePrefix) {
				return name
		}
		return captcha.DriverInstancePrefix + name
}

func (this *VerificationCodeService) Verify(id string, code string, ty string, clean ...bool) bool {
		var handler *captcha.Wrapper
		if len(clean) == 0 {
				clean = append(clean, true)
		}
		switch ty {
		case TypeMobile:
				handler = this.Mobile(id)
		default:
				handler = this.Image(ty)
		}
		return handler.Verify(id, code, clean...)
}

func (this *VerificationCodeService) SetStore(name string, instance base64Captcha.Store) {
		this.Stores[name] = instance
}

func (this *VerificationCodeService) SetDriver(name string, instance base64Captcha.Driver) {
		this.DriversContainer[name] = instance
}
