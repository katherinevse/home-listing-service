package config

type Config struct {
	PostgresCfg PostgresConfig `yaml:"db"`
}

type PostgresConfig struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DbName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func NewConfig() *Config {
	return &Config{
		PostgresCfg: PostgresConfig{},
	}
}
