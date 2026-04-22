package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	AryanApp      AryanApp      `yaml:"ARYAN_APP"`
	FararavandApp FararavandApp `yaml:"FARARAVAND_APP"`
	Database      Database      `yaml:"DATABASE"`
	App           App           `yaml:"APP"`
}

type AryanApp struct {
	BaseURL  string        `yaml:"BASE_URL"`
	APIKey   string        `yaml:"API_KEY"`
	UserName string        `yaml:"UserName"`
	Pass     string        `yaml:"Pass"`
	Timeout  time.Duration `yaml:"TIMEOUT"`
}

type FararavandApp struct {
	BaseURL  string        `yaml:"BASE_URL"`
	APIKey   string        `yaml:"API_KEY"`
	UserName string        `yaml:"USER_NAME"`
	Pass     string        `yaml:"PASSWORD"`
	Timeout  time.Duration `yaml:"TIMEOUT"`
}

type Database struct {
	Username           string `yaml:"USERNAME"`
	Password           string `yaml:"PASSWORD"`
	Port               string `yaml:"PORT"`
	Host               string `yaml:"HOST"`
	DBName             string `yaml:"DB_NAME"`
	MaxOpenConnections int    `yaml:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int    `yaml:"MAX_IDLE_CONNECTIONS"`
}

type App struct {
	LogPath string `yaml:"LOG_PATH"`
}

func Load(configPath string) (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetEnvPrefix("erp-job")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile(configPath)

	if err := v.MergeInConfig(); err != nil {
		return cfg, fmt.Errorf("read config %q: %w", configPath, err)
	}

	if err := v.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
