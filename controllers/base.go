package controllers

import (
		"encoding/json"
		"github.com/WebGameLinux/cms/dto/common"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/astaxie/beego"
)

type BaseController struct {
		beego.Controller
}

func (this *BaseController) GetJsonMapper() map[string]interface{} {
		var m = make(map[string]interface{})
		if err := json.Unmarshal(this.Ctx.Input.RequestBody, &m); err != nil {
				return nil
		}
		return m
}

func (this *BaseController) GetParamsMapper() mapper.Mapper {
		var m = make(map[string]interface{})
		if err := json.Unmarshal(this.Ctx.Input.RequestBody, &m); err != nil {
				return nil
		}
		return mapper.Mapper(m)
}

func (this *BaseController) ApiResponse(res common.JsonResponseInterface) {
		this.Data["json"] = res
		this.ServeJSON()
}
