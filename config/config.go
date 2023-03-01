package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	BaseURL string `yaml:"BASE_URL"`
	ApiKey  string `yaml:"API_KEY"`
	SqLite
}

type Server struct {
	Port string `yaml:"port"`
}

var Cfg config

func GetConfig() error {
	//open env file and read config
	file, err := os.Open("env.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil {
		return err
	}

	return nil
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