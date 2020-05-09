package common

import (
		"errors"
		"github.com/WebGameLinux/cms/dto/enums"
)

func Warn(code int, msg ...string) JsonResponseInterface {
		var res = NewJsonRepo(msg...)
		res.Set("error", errors.New("<Service:nil>"))
		res.Set("code", code)
		return res
}

func NewJsonRepo(msg ...string) JsonResponseInterface {
		var res = NewResult(true)
		if len(msg) == 0 {
				msg = append(msg, "ok")
		}
		res.Set("msg", msg[0])
		return res
}

func ParamError(code int, msg ...string) JsonResponseInterface {
		var res = NewJsonRepo(msg...)
		res.Set("error", errors.New("<Params:error>"))
		res.Set("code", code)
		return res
}

func ServerError(msg ...string) JsonResponseInterface {
		var res = NewJsonRepo(msg...)
		res.Set("error", errors.New("<Params:error>"))
		res.Set("code", enums.Error)
		return res
}

func Success(data interface{}, msg ...string) JsonResponseInterface {
		var res = NewJsonRepo(msg...)
		res.Set("data", data)
		res.Set("code", enums.SUCCESS)
		return res
}
