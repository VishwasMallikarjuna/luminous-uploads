package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
	DBName     string `yaml:"dbName"`
	DBHost     string `yaml:"dbHost"`
	DBPort     string `yaml:"dbPort"`
	ClientID   string `yaml:"clientID"`
	TenantID   string `yaml:"tenantID"`
	SecretKey  string `yaml:"secretKey"`
}

var AppConfig Config

func LoadConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "default/config/path/config.yml" // default path
	}
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}
