package config

type DefaultUserConfig struct {
	Avatars       []string `mapstructure:"avatars" yaml:"avatars"`
	ResetPassword string   `mapstructure:"reset_password" yaml:"reset_password"`
}
