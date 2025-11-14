package config

import "flag"

type AppFlags struct {
	ConfigPath string
}

func ParseFlags() AppFlags {
	configPath := flag.String("config", "./config/config.yaml", "path to config")
	flag.Parse()

	return AppFlags{
		ConfigPath: *configPath,
	}
}
