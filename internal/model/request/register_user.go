package request

type RegisterUser struct {
	Username string `json:"username" label:"用户名" binding:"required,min=3,max=32,alphanumunicode"`
	Nickname string `json:"nickname" label:"昵称" binding:"required,min=5,max=24"`
	Password string `json:"password" label:"密码" binding:"required,min=8,max=64,strongpassword"`
	Gender   int8   `json:"gender" label:"性别" binding:"oneof=1 2"`
}
