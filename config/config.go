package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	GeminiAPIKey string
}

var Cfg Config
var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		instance, err := loadConfig()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		Cfg = instance
	})
	return Cfg
}

func loadConfig() (Config, error) {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))

	return Config{
		Port:         os.Getenv("PORT"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}, err

}
