package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	AryanApp      AryanApp      `yaml:"ARYAN_APP"`
	FararavandApp FararavandApp `yaml:"FARARAVAND_APP"`
	Database      Database      `yaml:"DATABASE"`
	OTel          OTel          `yaml:"OTEL"`
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

type OTel struct {
	Enabled     bool   `yaml:"ENABLED"`
	Endpoint    string `yaml:"ENDPOINT"`
	Insecure    bool   `yaml:"INSECURE"`
	ServiceName string `yaml:"SERVICE_NAME"`
	Environment string `yaml:"ENVIRONMENT"`
}

type App struct {
	LogPath string `yaml:"LOG_PATH"`
}

func Load(configPath string) (Config, error) {
	return load(configPath, func(cfg Config) error {
		return cfg.Validate()
	})
}

func LoadDatabase(configPath string) (Database, error) {
	cfg, err := load(configPath, func(cfg Config) error {
		return cfg.ValidateDatabase()
	})
	if err != nil {
		return Database{}, err
	}

	return cfg.Database, nil
}

func load(configPath string, validate func(Config) error) (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetEnvPrefix("ERP_JOB")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile(configPath)
	v.SetDefault("OTEL.SERVICE_NAME", "erp-job")
	v.SetDefault("OTEL.ENVIRONMENT", "production")

	if err := v.MergeInConfig(); err != nil {
		return cfg, fmt.Errorf("read config %q: %w", configPath, err)
	}

	if err := v.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	if validate != nil {
		if err := validate(cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if err := c.ValidateDatabase(); err != nil {
		return err
	}

	if err := validateHTTPConfig("ARYAN_APP", c.AryanApp.BaseURL, c.AryanApp.APIKey, c.AryanApp.Timeout); err != nil {
		return err
	}

	if err := validateHTTPConfig("FARARAVAND_APP", c.FararavandApp.BaseURL, c.FararavandApp.APIKey, c.FararavandApp.Timeout); err != nil {
		return err
	}

	return c.validateOTel()
}

func (c Config) ValidateDatabase() error {
	if c.Database.Username == "" {
		return fmt.Errorf("DATABASE.USERNAME is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DATABASE.PASSWORD is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DATABASE.HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DATABASE.PORT is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DATABASE.DB_NAME is required")
	}
	if c.Database.MaxOpenConnections <= 0 {
		return fmt.Errorf("DATABASE.MAX_OPEN_CONNECTIONS must be greater than zero")
	}
	if c.Database.MaxIdleConnections < 0 {
		return fmt.Errorf("DATABASE.MAX_IDLE_CONNECTIONS cannot be negative")
	}
	if c.Database.MaxIdleConnections > c.Database.MaxOpenConnections {
		return fmt.Errorf("DATABASE.MAX_IDLE_CONNECTIONS cannot exceed DATABASE.MAX_OPEN_CONNECTIONS")
	}

	return nil
}

func (c Config) validateOTel() error {
	if c.OTel.Enabled {
		if c.OTel.Endpoint == "" {
			return fmt.Errorf("OTEL.ENDPOINT is required when OTEL.ENABLED is true")
		}
		if _, err := url.ParseRequestURI(c.OTel.Endpoint); err != nil {
			return fmt.Errorf("OTEL.ENDPOINT is invalid: %w", err)
		}
		if c.OTel.ServiceName == "" {
			return fmt.Errorf("OTEL.SERVICE_NAME is required when OTEL.ENABLED is true")
		}
		if c.OTel.Environment == "" {
			return fmt.Errorf("OTEL.ENVIRONMENT is required when OTEL.ENABLED is true")
		}
	}

	return nil
}

func validateHTTPConfig(prefix, rawURL, apiKey string, timeout time.Duration) error {
	if rawURL == "" {
		return fmt.Errorf("%s.BASE_URL is required", prefix)
	}
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return fmt.Errorf("%s.BASE_URL is invalid: %w", prefix, err)
	}
	if apiKey == "" {
		return fmt.Errorf("%s.API_KEY is required", prefix)
	}
	if timeout <= 0 {
		return fmt.Errorf("%s.TIMEOUT must be greater than zero", prefix)
	}

	return nil
}
