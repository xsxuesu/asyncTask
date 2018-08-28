package model

//自定义异常
type CustomException struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
