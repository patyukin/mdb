package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Logger struct {
		Level string `yaml:"level" validate:"required,oneof=debug info warn error dpanic panic fatal"`
		Mode  string `yaml:"mode" validate:"required,oneof=devel prod"`
	}
}

func LoadConfig() (*Config, error) {
	yamlConfigFilePath := os.Getenv("YAML_CONFIG_FILE_PATH")
	if yamlConfigFilePath == "" {
		return nil, fmt.Errorf("yaml config file path is not set")
	}

	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Printf("unable to close config file: %v", err)
		}
	}(f)

	var config Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
