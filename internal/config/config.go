package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Ollama  OllamaConfig  `yaml:"ollama"`
	Watcher WatcherConfig `yaml:"watcher"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type OllamaConfig struct {
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

type WatcherConfig struct {
	PollInterval int `yaml:"poll_interval"` // in milliseconds
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Ollama.Model == "" {
		cfg.Ollama.Model = "llama3"
	}
	if cfg.Watcher.PollInterval == 0 {
		cfg.Watcher.PollInterval = 1000
	}
	if cfg.Ollama.Host == "" {
		cfg.Ollama.Host = "http://localhost:11434"
	}

	return &cfg, nil
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
