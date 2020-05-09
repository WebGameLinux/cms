package request

type LoginParamsDto struct {
		Account  string `json:"account" validate:"required,gte=3,lte=100"` // 用户名, 邮箱 ，手机号
		Password string `json:"password" validate:"required,gt=6"`         // 密码
}

type CaptchaDto struct {
		Captcha string `json:"captcha"` // 验证码
}

type XsrfDto struct {
		XsrfToken string `json:"_xsrf"` // 表单安全验证
}

