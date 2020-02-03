package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host string
	Port int
	DB   DBConfig
}

type DBConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBname string
}

func ReadConfig(fname string) (*Config, error) {
	fullPath, _ := filepath.Abs(fname)
	yamlFile, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
