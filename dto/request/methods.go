package request

import (
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/validator"
		validation "gopkg.in/go-playground/validator.v9"
)

func (this *LoginParamsDto) Load(data map[string]interface{}) *LoginParamsDto {
		m := mapper.Mapper(data)
		if m.Len() == 0 {
				return this
		}
		this.Account = m.Get("account")
		this.Password = m.Get("password")
		return this
}

func (this *LoginParamsDto) Valid() (*validation.Validate, error) {
		var (
				err error
				v   = validator.GetValidator()
		)
		if err = v.Struct(this); err != nil {
				return v, err
		}
		return nil, nil
}