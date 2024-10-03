package config

type Config struct {
	AppCfg      AppConfig      `yaml:"app"`
	PostgresCfg PostgresConfig `yaml:"db"`
	JWTConfig   JWTConfig      `yaml:"jwt"`
}

type PostgresConfig struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DbName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type AppConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
	Domain string `yaml:"domain"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
}

func NewConfig() *Config {
	return &Config{
		PostgresCfg: PostgresConfig{},
		AppCfg:      AppConfig{},
		JWTConfig:   JWTConfig{},
	}
}
