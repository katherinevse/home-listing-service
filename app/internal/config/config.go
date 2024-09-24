package config

type Config struct {
	PostgresCfg PostgresConfig `yaml:"db"`
	//JWTConfig   JWTConfig      `yaml:"jwt"`
}

type PostgresConfig struct {
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DbName   string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

//type JWTConfig struct {
//	SecretKey string `yaml:"secret_key"`
//}

func NewConfig() *Config {
	return &Config{
		PostgresCfg: PostgresConfig{},
		//JWTConfig:   JWTConfig{},
	}
}
