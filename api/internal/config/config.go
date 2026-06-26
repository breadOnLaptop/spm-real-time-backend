package config

import "os"

type Config struct {
	DatabaseURL string
	RedisURL    string
	Port        string
}

func Load() *Config {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://spm_admin:secure_cryptographic_password@spm_db:5432/spm_telemetry?sslmode=disable"
	}
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "redis://spm_cache:6379/0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return &Config{
		DatabaseURL: dbUrl,
		RedisURL:    redisUrl,
		Port:        port,
	}
}
