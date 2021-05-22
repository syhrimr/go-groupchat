package model

type Config struct {
	DB    PostgreCfg `yaml:"postgresql"`
	Redis RedisCfg   `yaml:"redis"`
}

type PostgreCfg struct {
	Address  string `yaml:"address"`
	DBName   string `yaml:"db-name"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
}

type RedisCfg struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}
