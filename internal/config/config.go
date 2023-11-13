package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Database struct {
	URL string `yaml:"url"`
}

type NatsStreaming struct {
	URL       string `yaml:"url"`
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	Subject   string `yaml:"subject"`
}

type Server struct {
	URL string `yaml:"url"`
}

type Config struct {
	Database      Database      `yaml:"database"`
	NatsStreaming NatsStreaming `yaml:"nats-streaming"`
	Server        Server        `yaml:"http-server"`
}

var cfg *Config

func GetConfig() Config {
	if cfg != nil {
		return *cfg
	}
	return Config{}
}

func ReadConfigYAML(filepath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
