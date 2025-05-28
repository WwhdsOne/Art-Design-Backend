package config

type DefaultUserConfig struct {
	Avatars       []string `mapstructure:"avatars"`
	ResetPassword string   `mapstructure:"reset_password"`
}
