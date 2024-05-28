package config

type ProxyConfig struct {
	Scheme   string `mapstructure:"scheme"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
