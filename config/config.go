package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string

	// Database
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string

	// Oauth
	// Google
	GoogleClientID     string
	GoogleClientSecret string
	// Discord
	DiscordClientID     string
	DiscordClientSecret string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port: getEnv("PORT"),

		// Database
		DBHost:     getEnv("DB_HOST"),
		DBName:     getEnv("DB_NAME"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBPort:     getEnv("DB_PORT"),

		// Oauth
		// Google
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
		// Discord
		DiscordClientID:     getEnv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET"),
	}
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Printf("Env key not found: %s. Setting to an empty string.", key)
	return ""
}
