package validator

import (
		"github.com/go-playground/locales/zh"
		ut "github.com/go-playground/universal-translator"
		"gopkg.in/go-playground/validator.v9"
		zhTranslations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var trans ut.Translator

func GetTranslator() ut.Translator {
		if trans == nil {
				zh := zh.New()
				uni := ut.New(zh, zh)
				trans, _ = uni.GetTranslator("zh")
		}
		return trans
}

func GetValidator() *validator.Validate {
		validate := validator.New()
		_ = zhTranslations.RegisterDefaultTranslations(validate, GetTranslator())
		return validate
}

func TranslateError(err error, local ...string) string {
		if len(local) == 0 {
				local = append(local, "zh")
		}
		key := local[0]
		if errs, ok := err.(validator.ValidationErrors); ok {
				t := GetTranslator()
				info := errs.Translate(t)
				if v, ok := info[key]; ok {
						return v
				}
				for _, v := range info {
						return v
				}
		}
		return err.Error()
}
