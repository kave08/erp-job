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
	AryanApp      AryanApp      `yaml:"ARYAN_APP"`
	FararavandApp FararavandApp `yaml:"FARARAVAND_APP"`
	Database      Database      `yaml:"DATABASE"`
}
type AryanApp struct {
	BaseURL  string `yaml:"BASE_URL"`
	APIKey   string `yaml:"API_KEY"`
	UserName string `yaml:"UserName"`
	Pass     string `yaml:"Pass"`
}
type FararavandApp struct {
	BaseURL  string `yaml:"BASE_URL"`
	APIKey   string `yaml:"API_KEY"`
	UserName string `yaml:"USER_NAME"`
	Pass     string `yaml:"PASSWORD"`
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

type SetupResult struct {
	AryanApp        AryanApp
	MysqlConnection *sql.DB
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
		config.TagName = "yaml"
	})
	if err != nil {
		fmt.Println("Error in unmarshaling config")
		panic(err)
	}

	fmt.Printf("%v", Cfg)

	mdb, err := initializeMySQL(Cfg.Database)
	if err != nil {
		panic(fmt.Sprintf("error at connecting to mysql database. err: %v, connection info: %+v", err, Cfg.Database))
	}

	return &SetupResult{
		AryanApp:        Cfg.AryanApp,
		MysqlConnection: mdb,
	}
}
