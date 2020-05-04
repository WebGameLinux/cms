package drivers

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/captcha/items"
		"github.com/mojocn/base64Captcha"
		"time"
)

type CaptchaMobileDriver struct {
		Mobile string
		Length int
		Source string
}

// id, content, answer := c.Driver.GenerateIdQuestionAnswer()
//	item, err := c.Driver.DrawCaptcha(content)
//	if err != nil {
//		return "", "", err
//	}
//	c.Store.Set(id, answer)
//	b64s = item.EncodeB64string()
//	return

func NewCaptchaMobileDriver(mobile string, length int, source string) *CaptchaMobileDriver {
		var driver = new(CaptchaMobileDriver)
		driver.Mobile = mobile
		driver.Length = length
		driver.Source = source
		return driver
}

func (this *CaptchaMobileDriver) DrawCaptcha(content string) (item base64Captcha.Item, err error) {
		var it = new(items.MobileCodeStrItem)
		it.Code = content
		return it, nil
}

func (this *CaptchaMobileDriver) GenerateIdQuestionAnswer() (id, q, a string) {
		var an = this.answer()
		return this.id(), an, an
}

func (this *CaptchaMobileDriver) answer() string {
		return base64Captcha.RandText(this.length(), this.source())
}

func (this *CaptchaMobileDriver) id() string {
		return fmt.Sprintf("%s%d", base64Captcha.RandText(10, this.source()), time.Now().Unix())
}

func (this *CaptchaMobileDriver) length() int {
		if this.Length <= 0 || this.Length > 10 {
				return 6
		}
		return this.Length
}

func (this *CaptchaMobileDriver) source() string {
		if this.Source == "" {
				this.Source = base64Captcha.TxtNumbers
		}
		return this.Source
}

func RandText(size int) string {
		return base64Captcha.RandText(size, base64Captcha.TxtAlphabet+base64Captcha.TxtNumbers)
}
