package captcha

import (
		"github.com/mojocn/base64Captcha"
)

type Wrapper struct {
		Instance *base64Captcha.Captcha
		Config   map[string]interface{}
}

func NewCaptcha(config ...map[string]interface{}) *Wrapper {
		if len(config) == 0 {
				config = append(config, map[string]interface{}{})
		}
		var w = new(Wrapper)
		w.Config = config[0]

		return w
}

func (this *Wrapper) init() {
		conf := NewConfiguration(this.Config)
		this.Instance = base64Captcha.NewCaptcha(conf.GetDriverInstance(), conf.GetStoreInstance())
}

func (this *Wrapper) GetCaptcha() *base64Captcha.Captcha {
		if this.Instance == nil {
				this.init()
		}
		return this.Instance
}

func (this *Wrapper) Generate() (id, b64s string, err error) {
		return this.GetCaptcha().Generate()
}

func (this *Wrapper) Verify(id, text string, clean ...bool) bool {
		if len(clean) == 0 {
				clean = append(clean, true)
		}
		return this.GetCaptcha().Verify(id, text, clean[0])
}
