package config

import "time"

type HealthzConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Enabled      bool
}

var defaultHealthzConfig = HealthzConfig{
	Port:         "60005",
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
	Enabled:      true,
}
