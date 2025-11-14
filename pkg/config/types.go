package config

import "time"

type HTTPConfig struct {
	Address      string        `yaml:"address" env:"HTTP_ADDRESS" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"5s"`
}
