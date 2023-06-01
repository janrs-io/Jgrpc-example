package config

import (
	"github.com/spf13/viper"
)

// Config Service config
type Config struct {
	Grpc      Grpc      `json:"grpc" yaml:"grpc"`
	Http      Http      `json:"http" yaml:"http"`
	Redis     Redis     `json:"redis" yaml:"redis"`
	WhiteList WhiteList `json:"whiteList" yaml:"whiteList"`
	Trace     Trace     `json:"trace" yaml:"trace"`
}

// NewConfig Initial service's config
func NewConfig(cfg string) *Config {

	if cfg == "" {
		panic("load config file failed.config file can not be empty.")
	}

	viper.SetConfigFile(cfg)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		panic("read config failed.[ERROR]=>" + err.Error())
	}
	conf := &Config{}
	// Assign the overloaded configuration to the global
	if err := viper.Unmarshal(conf); err != nil {
		panic("assign config failed.[ERROR]=>" + err.Error())
	}

	return conf

}
