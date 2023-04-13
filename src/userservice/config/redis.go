package config

import "time"

// Redis Redis Config
type Redis struct {
	Host         string        `json:"host" yaml:"host"`
	Port         string        `json:"port" yaml:"port"`
	Username     string        `json:"username" yaml:"username"`
	Password     string        `json:"password" yaml:"password"`
	Database     int           `json:"database" yaml:"database"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout" yaml:"pool_timeout"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
}
