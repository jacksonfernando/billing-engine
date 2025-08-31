package global

type Configuration struct {
	HostUrl                      string `mapstructure:"host_url"`
	HostPort                     string `mapstructure:"host_port"`
	DbHost                       string `mapstructure:"db_host"`
	DbName                       string `mapstructure:"db_name"`
	DbUser                       string `mapstructure:"db_user"`
	DbPass                       string `mapstructure:"db_pass"`
	DbPort                       string `mapstructure:"db_port"`
	TimeoutDuration              int    `mapstructure:"timeout_duration"`
	PrivateJWTAccessTokenSecret  string `mapstructure:"private_jwt_access_token_secret"`
	PrivateJWTRefreshTokenSecret string `mapstructure:"private_jwt_refresh_token_secret"`
}
