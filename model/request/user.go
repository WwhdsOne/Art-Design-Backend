package request

type User struct {
	Username     string   `json:"username" label:"用户名" validate:"required,min=8,max=32,alphanumunicode"`
	RealName     string   `json:"real_name" label:"真实姓名" validate:"max=50"`
	Nickname     string   `json:"nickname" label:"昵称" validate:"required,min=5,max=24"`
	Password     string   `json:"password" label:"密码" validate:"required,min=8,max=64,strongpassword"`
	Gender       int8     `json:"gender" label:"性别" validate:"oneof=1 2"`
	Email        string   `json:"email" label:"邮箱" validate:"omitempty,email,max=100"`
	Phone        string   `json:"phone" label:"手机号" validate:"omitempty,e164,max=30"`
	Address      string   `json:"address" label:"地址" validate:"max=256"`
	Introduction string   `json:"introduction" label:"个人介绍" validate:"max=256"`
	Occupation   string   `json:"occupation" label:"职业" validate:"max=50"`
	Tags         []string `json:"tags" label:"个人标签" validate:"dive,max=20"`
}
