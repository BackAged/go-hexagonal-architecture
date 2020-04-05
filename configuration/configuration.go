package configuration

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Application holds application configurations
type Application struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Mongo           Mongo  `yaml:"mongo"`
	GracefulTimeout int    `yaml:"graceful_timeout"`
}

// Mongo holds mongodb configurations
type Mongo struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

// Load loads config from path
func Load(path string) (*Application, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer f.Close()

	cfg := Application{}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
