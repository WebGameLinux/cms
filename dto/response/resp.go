package response

type RespJson struct {
		Data interface{} `json:"data"`
		Msg  string      `json:"message"`
		Code int         `json:"code"`
}
