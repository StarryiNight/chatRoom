package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

// ParamLogin 登陆请求参数
type ParamLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
