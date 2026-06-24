// Package config loads app configuration from a YAML file + env overrides.
package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

// Config is the root configuration structure.
type Config struct {
	App   AppConfig   `yaml:"app"`
	Mongo MongoConfig `yaml:"mongo"`
	Redis RedisConfig `yaml:"redis"`
	JWT   JWTConfig   `yaml:"jwt"`
}

type AppConfig struct {
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

type MongoConfig struct {
	URI string `yaml:"uri"`
	DB  string `yaml:"db"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expiry string `yaml:"expiry"`
}

// Load reads the YAML file at path and returns a Config.
func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}

	// Allow environment overrides for sensitive values.
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.Mongo.URI = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}

	return &cfg, nil
}
