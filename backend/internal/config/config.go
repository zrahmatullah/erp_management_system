package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
	Expire int // hours
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env_database")
	_ = godotenv.Load(".env_database")

	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	dsn := getEnv("DATABASE_URL", "")
	if dsn == "" {
		dbUser := getEnv("DB_USERNAME", "postgres")
		dbPass := getEnv("DB_PASSWORD", "samp3321")
		dbHost := getEnv("DB_HOST", "127.0.0.1")
		dbPort := getEnv("DB_PORT", "5432")
		dbName := getEnv("DB_DATABASE", "cafe_erp")
		dsn = "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		DB: DBConfig{
			DSN: dsn,
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "127.0.0.1:6379"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "super-secret-cafe-erp-jwt-key-2026"),
			Expire: jwtExpire,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
