/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	v "github.com/spf13/viper"

	"gopkg.in/yaml.v2"
)

const (
	GPUCONF_ENV_PREFIX = "GPUCONF"
)

type Config struct {
	Viper *v.Viper `yaml:"-" json:"-"`

	General CGeneral `mapstructure:"general" json:"general,omitempty" yaml:"general,omitempty"`
	Logging CLogging `mapstructure:"logging" json:"logging,omitempty" yaml:"logging,omitempty"`
}

type CGeneral struct {
	Debug bool `mapstructure:"debug,omitempty" json:"debug,omitempty" yaml:"debug,omitempty"`
}

type CLogging struct {
	// Path of the logfile
	Path string `mapstructure:"path,omitempty" json:"path,omitempty" yaml:"path,omitempty"`
	// Enable/Disable logging to file
	EnableLogFile bool `mapstructure:"enable_logfile,omitempty" json:"enable_logfile,omitempty" yaml:"enable_logfile,omitempty"`
	// Enable JSON format logging in file
	JsonFormat bool `mapstructure:"json_format,omitempty" json:"json_format,omitempty" yaml:"json_format,omitempty"`

	// Log level
	Level string `mapstructure:"level,omitempty" json:"level,omitempty" yaml:"level,omitempty"`

	// Enable emoji
	EnableEmoji bool `mapstructure:"enable_emoji,omitempty" json:"enable_emoji,omitempty" yaml:"enable_emoji,omitempty"`
	// Enable/Disable color in logging
	Color bool `mapstructure:"color,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
}

func NewConfig(viper *v.Viper) *Config {
	if viper == nil {
		viper = v.New()
	}

	GenDefault(viper)
	return &Config{Viper: viper}
}

func (c *Config) GetGeneral() *CGeneral {
	return &c.General
}

func (c *Config) GetLogging() *CLogging {
	return &c.Logging
}

func (c *Config) Unmarshal() error {
	c.Viper.ReadInConfig()

	err := c.Viper.Unmarshal(&c)

	return err
}

func (c *Config) Yaml() ([]byte, error) {
	return yaml.Marshal(c)
}

func GenDefault(viper *v.Viper) {
	viper.SetDefault("general.debug", false)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.enable_logfile", false)
	viper.SetDefault("logging.path", "/var/log/gpu-configurator.log")
	viper.SetDefault("logging.json_format", false)
	viper.SetDefault("logging.enable_emoji", true)
	viper.SetDefault("logging.color", true)
}

func (g *CGeneral) HasDebug() bool {
	return g.Debug
}
