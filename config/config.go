package config

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var Cfg config

type config struct {
	BaseURL string `yaml:"BASE_URL"`
	ApiKey  string `yaml:"API_KEY"`
	SqLite  SqLite `yaml:"SQLITE"`
}

type Server struct {
	Port string `yaml:"port"`
}

type SqLite struct {
	Username           string `yml:"Username"`
	Password           string `yml:"Password"`
	Host               string `yml:"Host"`
	Port               string `yml:"Port"`
	DBName             string `yml:"DB_Name"`
	MaxOpenConnections int    `yml:"Max_Open_Connections"`
	MaxIdleConnections int    `yml:"Max_Idle_Connections"`
}

type SetupResult struct {
	SqlitConnection *sql.DB
}

func LoadConfig(configPath string) *SetupResult {

	viper.SetEnvPrefix("erp-job")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.MergeInConfig()
	if err != nil {
		fmt.Println("Error in reading config")
		panic(err)
	}

	err = viper.Unmarshal(&Cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yml"
	})
	if err != nil {
		fmt.Println("Error in unmarshaling config")
		panic(err)
	}

	fmt.Printf("%v", Cfg)

	sdb, err := initializeSQLite()
	if err != nil {
		panic(err)
	}

	return &SetupResult{
		SqlitConnection: sdb,
	}
}
