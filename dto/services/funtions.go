package services

import "errors"

func NewEmptyResult(code int, msg ...string) ResultStruct {
		var res = NewResultByMsg(msg...)
		res.Set("error", errors.New("<Service:nil>"))
		res.Set("code", code)
		return res
}

func NewResultByMsg(msg ...string) ResultStruct {
		var res = NewResult(true)
		if len(msg) == 0 {
				msg = append(msg, "ok")
		}
		res.Set("msg", msg[0])
		return res
}

func NewParamsErrorResult(code int, msg ...string) ResultStruct {
		var res = NewResultByMsg(msg...)
		res.Set("error", errors.New("<Params:error>"))
		res.Set("code", code)
		return res
}

func NewServerErrorResult(msg ...string) ResultStruct {
		var res = NewResultByMsg(msg...)
		res.Set("error", errors.New("<Params:error>"))
		res.Set("code", 5000)
		return res
}

func NewSuccessResult(data interface{}, msg ...string) ResultStruct {
		var res = NewResultByMsg(msg...)
		res.Set("data", data)
		res.Set("code", 0)
		return res
}
