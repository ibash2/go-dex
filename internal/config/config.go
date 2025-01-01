package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	// db "url-shortener/internal/storage/postgre"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	// StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`

	// DataBase db.DB `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8084"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// Указываем полный путь к конфигурационному файлу
	configPath := filepath.Join("/home/abdu1bari/go/projects/url-shortener/config/local.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
