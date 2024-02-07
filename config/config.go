package config

import (
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type Configs struct {
	App   *AppConfig
	DB    *DBConf
	Redis *RedisConf
	Token *Token
}

type AppConfig struct {
	TimeOut time.Duration
	Port    int
}

type DBConf struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
	TimeOut  time.Duration
}

type RedisConf struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}
type Token struct {
	Refresh *TokenConf
	Access  *TokenConf
}
type TokenConf struct {
	TokenSecret string
	ExpiresAt   time.Duration
}

func New() (*Configs, error) {
	configFile := "config.yaml"
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Configs{}

	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
