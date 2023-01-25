package config

import (
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server `yaml:"server"`
}

type Server struct {
	Port int `yaml:"port"`
}

var LoadConfig Config

func GetConfig() error {
	//open env file and read config
	file, err := os.Open("env.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&LoadConfig)
	if err != nil {
		return err
	}

	return nil
}