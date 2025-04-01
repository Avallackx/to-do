package bootstrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
}

func LoadEnv() Env {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	env := Env{
		AppEnv:        getEnv("APP_ENV", "development"),
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "admin"),
		DBName:        getEnv("DB_NAME", "todo"),
	}

	return env
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
