package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	AppEnv              string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	ServerPort          string
	JwtSecret           string
	RedisAddr           string
	CloudinaryCloudName string
	CloudinaryApiKey    string
	CloudinaryApiSecret string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %q is not set", key)
	}
	return v
}

func Load() *Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	if appEnv != "production" {
		godotenv.Load(".env." + appEnv)
		godotenv.Load()
	}

	return &Config{
		AppEnv:              appEnv,
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "postgres"),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		JwtSecret:           mustGet("JWT_SECRET"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", "test"),
		CloudinaryApiKey:    getEnv("CLOUDINARY_API_KEY", "test"),
		CloudinaryApiSecret: getEnv("CLOUDINARY_API_SECRET", "test"),
	}
}

func (c *Config) Validate() error {
	if c.DBPassword == "" && c.IsProd() {
		return fmt.Errorf("DB_PASSWORD must be set in production")
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func (c *Config) RedisOptions() *redis.Options {
	return &redis.Options{
		Addr: c.RedisAddr,
	}
}

func (c *Config) IsProd() bool { return c.AppEnv == "production" }
func (c *Config) IsDev() bool  { return c.AppEnv == "dev" }
func (c *Config) IsTest() bool { return c.AppEnv == "test" }
