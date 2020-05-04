package captcha

import (
		"github.com/mojocn/base64Captcha"
		"image/color"
		"math/rand"
)

func (this *ConfigurationWrapper) GetId() string {
		return this.Id
}

func (this *ConfigurationWrapper) GetType() string {
		return this.CaptchaType
}

func (this *ConfigurationWrapper) GetDriver() string {
		if this.CaptchaType == "" {
				this.CaptchaType = DriverAudio
		}
		return this.CaptchaType
}

func (this *ConfigurationWrapper) GetDriverInstance() base64Captcha.Driver {
		if this.DriverInstance == nil {
				this.DriverInstance = GetDriver(this.Driver)
		}
		return this.DriverInstance
}

func (this *ConfigurationWrapper) GetStoreInstance() base64Captcha.Store {
		if this.StoreInstance == nil {
				return base64Captcha.DefaultMemStore
		}
		return this.StoreInstance
}

// 获取启动
func GetDriver(name string) base64Captcha.Driver {
		var (
				height          = 80
				width           = 240
				noiseCount      = 5
				length          = 6
				source          string
				fonts           []string
				showLineOptions = rand.Intn(3)
				bgColor         = &color.RGBA{R: 0, G: 0, B: 0, A: 0}
		)

		switch name {
		case DriverAudio:
				return base64Captcha.NewDriverAudio(length, "en")
		case DriverString:
				source = base64Captcha.TxtChineseCharaters
				source += base64Captcha.TxtAlphabet
				source += base64Captcha.TxtNumbers
				return base64Captcha.NewDriverString(height, width, noiseCount, showLineOptions, length, source, bgColor, fonts).ConvertFonts()
		case DriverMath:
				return base64Captcha.NewDriverMath(height, width, noiseCount, showLineOptions, bgColor, fonts).ConvertFonts()
		case DriverChinese:
				source = base64Captcha.TxtChineseCharaters
				return NewDriverChinese(height, width, noiseCount, showLineOptions, length, source, bgColor).ConvertFonts()
		case DriverEnglish:
				source = base64Captcha.TxtAlphabet
				return base64Captcha.NewDriverLanguage(height, width, noiseCount, showLineOptions, length, bgColor, nil, source)
		case DriverDigit:
				fallthrough
		default:
				return base64Captcha.DefaultDriverDigit
		}
}

// 创建一个中文去驱动
func NewDriverChinese(height int, width int, noiseCount int, showLineOptions int, length int, source string, bgColor *color.RGBA) *base64Captcha.DriverChinese {
		return &base64Captcha.DriverChinese{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, Source: source, BgColor: bgColor}
}
