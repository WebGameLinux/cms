package middlewares

import (
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/response"
		"github.com/WebGameLinux/cms/dto/services"
		"github.com/WebGameLinux/cms/models"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"log"
		"strconv"
)

func Auth() beego.FilterFunc {
		return token
}

func token(ctx *context.Context) {
		token := ctx.Request.Header.Get(enums.AuthToken)
		if token == "" {
				c, err := ctx.Request.Cookie(enums.AuthToken)
				if err != nil {
						UnLogin(ctx)
						return
				}
				token = c.Value
		}
		if !services.GetTokenService().Verify(token) {
				TokenFail(ctx)
				return
		}
		user := new(models.User)
		if !services.GetTokenService().Load(token, user) {
				TokenError(ctx)
				return
		}
		//	services.GetTokenService().Dispatch(token,user)
		ctx.Request.Header.Set("uid", strconv.Itoa(int(user.Id)))
		return
}

func UnLogin(ctx *context.Context) {
		errorInfo := response.RespJson{
				Data: nil,
				Msg:  enums.ErrorUserUnLogin.WrapMsg(),
				Code: enums.ErrorUserUnLogin.Int(),
		}
		info, _ := errorInfo.MarshalJSON()
		ResponseApi(ctx, info)
		return
}

func TokenFail(ctx *context.Context) {
		errorInfo := response.RespJson{
				Data: nil,
				Msg:  enums.ErrorAuthCheckTokenFail.WrapMsg(),
				Code: enums.ErrorAuthCheckTokenFail.Int(),
		}
		info, _ := errorInfo.MarshalJSON()
		ResponseApi(ctx, info)
		return
}

func TokenError(ctx *context.Context) {
		errorInfo := response.RespJson{
				Data: nil,
				Msg:  enums.ErrorAuth.WrapMsg(),
				Code: enums.ErrorAuth.Int(),
		}
		info, _ := errorInfo.MarshalJSON()
		ResponseApi(ctx, info)
		return
}

func ResponseApi(ctx *context.Context, info []byte) {
		ctx.Output.Header("Content-Type", "application/json; charset=UTF-8")
		if _, err := ctx.ResponseWriter.Write(info); err != nil {
				log.Fatal(err)
		}
}
