package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func loadEnvStr(key string, res *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	*res = s
}

func loadEnvUint(key string, res *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	si, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	*res = uint(si)
}

type config struct {
	Listen   listenConfig `yaml:"listen" json:"listen"`
	DBCfg    pgConfig     `yaml:"db" json:"db"`
	TokenCfg tokenConfig  `yaml:"token" json:"token"`
	RedisCfg redisConfig  `yaml:"redis" json:"redis"`
}

func (c *config) loadFromEnv() {
	c.Listen.loadFromEnv()
	c.DBCfg.loadFromEnv()
	c.TokenCfg.loadFromEnv()
	c.RedisCfg.loadFromEnv()
}

func defaultConfig() config {
	return config{
		Listen:   defaultListenConfig(),
		DBCfg:    defaultPgConfig(),
		TokenCfg: defaultTokenConfig(),
		RedisCfg: defaultredisConfig(),
	}
}

func loadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)
	if err != nil {
		return err
	}

	fl, err := os.Open(filepath.Clean(fn))
	if err != nil {

		return err
	}

	defer fl.Close()
	return yaml.NewDecoder(fl).Decode(c)
}

func LoadConfig(fn string) config {
	cfg := defaultConfig()
	cfg.loadFromEnv()
	if err := loadConfigFromFile(fn, &cfg); err != nil {
		if err != nil {
			log.Warn().Str("file", fn).Err(err).Msg("cannot load config file, use defaults")
		}
	}
	return cfg
}
