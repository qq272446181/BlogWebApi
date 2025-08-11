package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	ExpireHours int    `json:"expire_hours"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
}

var (
	config *AppConfig
	once   sync.Once
)

// LoadConfig 从JSON文件加载配置
func LoadConfig() *AppConfig {
	once.Do(func() {
		file, err := os.Open("./config.json")
		if err != nil {
			log.Fatalf("无法打开配置文件: %v", err)
		}
		defer file.Close()

		config = &AppConfig{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			log.Fatalf("无法解析配置文件: %v", err)
		}
	})
	return config
}

// GetConfig 获取应用配置
func GetConfig() *AppConfig {
	if config == nil {
		return LoadConfig()
	}
	return config
}
