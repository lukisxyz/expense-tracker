package config

type tokenConfig struct {
	Key             string `yaml:"key" json:"key"`
	AccessDuration  uint   `yaml:"access_duration" json:"access_duration"`
	RefreshDuration uint   `yaml:"refresh_duration" json:"refresh_duration"`
}

func defaultTokenConfig() tokenConfig {
	return tokenConfig{
		Key:             "12345678901234567890123456789012",
		AccessDuration:  15,
		RefreshDuration: 30,
	}
}

func (l *tokenConfig) loadFromEnv() {
	loadEnvStr("JWT_TOKEN", &l.Key)
	loadEnvUint("ACCESS_TOKEN_DURATION", &l.AccessDuration)
	loadEnvUint("REFRESH_TOKEN_DURATION", &l.RefreshDuration)
}
