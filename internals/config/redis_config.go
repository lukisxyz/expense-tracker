package config

import "fmt"

type redisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
}

func (p *redisConfig) ConnStr() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func defaultredisConfig() redisConfig {
	return redisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "password",
	}
}

func (p *redisConfig) loadFromEnv() {
	loadEnvStr("REDIS_HOST", &p.Host)
	loadEnvUint("REDIS_PORT", &p.Port)
	loadEnvStr("REDIS_PASSWORD", &p.Password)
}
