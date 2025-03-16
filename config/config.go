package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance *Config
	once     sync.Once
)

const (
	DEFAULT_HOST        = "0.0.0.0"
	DEFAULT_PORT        = 8080
	DEFAULT_CONFIG_PATH = "config/config.yaml"
	DEFAULT_STUN_SERVER = "stun:stun.l.google.com:19302"
)

type Http struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type WebRTC struct {
	StunServer string `yaml:"stun_server"`
}

type Config struct {
	Http   Http   `yaml:"http"`
	WebRTC WebRTC `yaml:"webrtc"`
}

func GetDefaultConfig() *Config {
	return &Config{
		Http: Http{
			Host: DEFAULT_HOST,
			Port: DEFAULT_PORT,
		},
		WebRTC: WebRTC{
			StunServer: DEFAULT_STUN_SERVER,
		},
	}
}

// GetConfig 返回 Config 的单例实例
func GetConfig() *Config {
	once.Do(func() {
		instance = NewConfig()
	})
	return instance
}

func NewConfig() *Config {
	var cfg = GetDefaultConfig()

	data, err := os.ReadFile(DEFAULT_CONFIG_PATH)
	if err != nil {
		log.Printf("Failed to read config file %s: %v, using default config", DEFAULT_CONFIG_PATH, err)
		return cfg
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Failed to parse config file: %v, using default config", err)
		return cfg
	}

	return cfg
}
