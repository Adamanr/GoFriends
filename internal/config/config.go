package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Salt     string `yaml:"salt"`
	} `yaml:"database"`
	CS struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"chi-server"`
	Minio struct {
		AccessKey string `yaml:"access-key"`
		SecretKey string `yaml:"secret-key"`
		Endpoint  string `yaml:"endpoint"`
		SSL       bool   `yaml:"ssl"`
	} `yaml: "minio-db"`
}

func GetConfigs() (*Config, error) {
	file, err := os.ReadFile("./locale.yaml")
	if err != nil {
		return nil, err
	}

	var cfg *Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
