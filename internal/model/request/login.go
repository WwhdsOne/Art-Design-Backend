package request

type Login struct {
	Username string `json:"username" label:"用户名" binding:"required"`
	Password string `json:"password" label:"密码" binding:"required"`
}
