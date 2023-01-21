package config

import "os"

type Config struct {
	Server `yaml:"server"`
}

type Server struct {
	Port int `yaml:"port"`
}

func GetConfig() error {
	//open env file and read config
	file, err := os.Open("/env.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return nil
}