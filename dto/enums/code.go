package enums

type Code int

const (
		SUCCESS          Code = 200
		ERROR            Code = 500
		InvalidParams    Code = 400
		MethodNotAllowed Code = 405
		FORBIDDEN        Code = 403

		ErrorFileIsEmpty     Code = 10001
		ErrorFileIsTooLarge  Code = 10002
		ErrorCanNotGetImgUrl Code = 10003
		ErrorTooManyImages   Code = 10004
		ErrorFileType        Code = 10005

		ErrorAuthCheckTokenFail    Code = 20001
		ErrorAuthCheckTokenTimeout Code = 20002
		ErrorAuthToken             Code = 20003
		ErrorAuth                  Code = 20004
		ErrorAccessDenied          Code = 20005

		ErrorTaskRepeat Code = 30001

		UpdateModelField Code = 3100 // 更新失败
		RecordNotExists  Code = 3101 // 数据记录不存在
		DataServiceError Code = 3102 // 数据查询服务异常
		CreateRecordField Code = 3103 // 创建记录失败

		//USER ERROR 4xxx

		ErrorUserLogin  Code = 40001
		ErrorUserCookie Code = 40002

		ErrorUserAlreadyExist Code = 40003
		ErrorUserNotExist     Code = 40004
		ErrorUserResetToken   Code = 40005
		ErrorUserUnLogin      Code = 40006
		ErrorUserDelete       Code = 40007

		//上传参数错误(api 选择)
		ErrorUploadParam      Code = 50001
		ErrorCanNotUpload     Code = 50002
		ErrorUploadTokenError Code = 50003
		// 未定义业务码
		UnknownCode  Code = 55000
		UnknownError Code = 55001
)
