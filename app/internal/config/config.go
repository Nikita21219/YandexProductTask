package config

type PostgresConfig struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DbName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Config struct {
	Host        string         `yaml:"host"`
	Port        string         `yaml:"port"`
	PostgresCfg PostgresConfig `yaml:"db"`
}

func NewConfig() *Config {
	return &Config{
		Host:        "",
		Port:        "",
		PostgresCfg: PostgresConfig{},
	}
}
