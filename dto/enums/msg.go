package enums

import (
		"fmt"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/reflects"
)

var msgFlags = map[Code]string{
		SUCCESS:                    "ok",
		ERROR:                      "fail",
		InvalidParams:              "请求参数错误",
		ErrorFileIsEmpty:           "上传文件为空",
		ErrorFileIsTooLarge:        "上传文件太大",
		ErrorFileType:              "文件类型错误",
		ErrorCanNotGetImgUrl:       "无法获取第三方图床 URL",
		ErrorTooManyImages:         "上传图片太多",
		ErrorAuthCheckTokenFail:    "Token鉴权失败",
		ErrorAuthCheckTokenTimeout: "Token 过期",
		ErrorAuthToken:             "Token 不正确",
		ErrorAuth:                  "认证失败",
		ErrorAccessDenied:          "禁止访问",

		ErrorTaskRepeat:   "任务重复提交,请等待当前任务完成",
		UpdateModelField:  "更新失败",
		RecordNotExists:   "数据记录不存在",
		DataServiceError:  "数据查询服务异常",
		CreateRecordField: "创建记录失败",
		ParamEmpty:        "空参数异常",
		KeyExists:         "记录已存在",
		CreateTokenFailed: "登陆异常",
		ErrorUserLogin:    "用户不存在或用户名密码错误",

		ErrorUserCookie: "用户 COOKIE 错误",

		ErrorUserAlreadyExist: "用户已存在",
		ErrorUserNotExist:     "用户不存在",
		ErrorUserResetToken:   "重置 Token 错误",
		ErrorUserUnLogin:      "用户未登录",
		Error:                 "服务异常",
		ErrorUploadParam:      "上传参数错误",
		ErrorCanNotUpload:     "无法上传图片到第三方图床",
		ErrorUploadTokenError: "上传 Token 错误",

		UnknownCode:  "未定义业务码",
		UnknownError: "未知异常",
}

func GetMsg(code int) string {
		msg, ok := msgFlags[Code(code)]
		if ok {
				return msg
		}
		return msgFlags[UnknownCode]
}

func Parse(msg string) Code {
		for code, m := range msgFlags {
				if m == msg {
						return code
				}
		}
		return UnknownCode
}

func GetMsgList() map[Code]string {
		return msgFlags
}

func String(code Code) string {
		return code.String()
}

func (this Code) String() string {
		return GetMsg(int(this))
}

func (this Code) Parse(msg string) Code {
		return Parse(msg)
}

func (this Code) Json() string {
		return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, int(this), this.String())
}

func (this Code) Map() map[Code]string {
		var m = make(map[Code]string)
		m[this] = this.String()
		return m
}

func (this Code) Int() int {
		return int(this)
}

func (this Code) WrapMsg(err ...interface{}) string {
		if len(err) == 0 {
				return this.String()
		}
		wrap := err[0]
		if wrap == nil {
				return this.String()
		}
		if e, ok := wrap.(error); ok {
				return this.String() + "," + e.Error()
		}
		if str, ok := wrap.(string); ok && str != "" {
				return this.String() + "," + str
		}
		return this.String() + "\n" + reflects.Any2Str(wrap)
}

func (this Code) Replace(err ...interface{}) string {
		if len(err) == 0 {
				return this.String()
		}
		wrap := err[0]
		if wrap == nil {
				return this.String()
		}
		if e, ok := wrap.(error); ok {
				return e.Error()
		}
		if str, ok := wrap.(string); ok && str != "" {
				return str
		}
		if str, ok := wrap.(fmt.Stringer); ok {
				return str.String()
		}
		return this.String()
}

func (this Code) Struct() *models.BusinessCode {
		var code = new(models.BusinessCode)
		code.Code = this.Int()
		code.Message = this.String()
		return code
}

func (this Code) Equal(code Code) bool {
		return this == code
}
