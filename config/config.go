package config

import (
	"github.com/spf13/viper"
)

const (
	PATH     = "PATH"
	HttpPort = "HttpPort"
)

type (
	Provider interface {
		Config() *Config
	}
	Config struct{}
)

func (p *Config) PATH() string {
	return viper.GetString(PATH)
}

func (p *Config) HttpPort() string {
	return viper.GetString(HttpPort)
}
