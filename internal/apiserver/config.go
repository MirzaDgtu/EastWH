package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL: "pmp:pmp1226@(nor.ru:3306)/eastwh?parseTime=true",
		BindAddr:    "127.0.0.1:8091",
		LogLevel:    "debug",
	}
}
