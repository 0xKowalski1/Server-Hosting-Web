package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	// Database
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string

	// Session Store
	SessionStoreSecret   string
	SessionStoreHttpOnly bool
	SessionStoreSecure   bool

	// Oauth
	// Google
	GoogleClientID     string
	GoogleClientSecret string
	// Discord
	DiscordClientID     string
	DiscordClientSecret string

	// Stripe
	StripePublicKey string
	StripeSecretKey string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port: getEnvString("PORT"),

		// Database
		DBHost:     getEnvString("DB_HOST"),
		DBName:     getEnvString("DB_NAME"),
		DBUser:     getEnvString("DB_USER"),
		DBPassword: getEnvString("DB_PASSWORD"),
		DBPort:     getEnvString("DB_PORT"),

		// Session Store
		SessionStoreSecret:   getEnvString("SESSION_STORE_SECRET"),
		SessionStoreHttpOnly: getEnvBool("SESSION_STORE_HTTP_ONLY"),
		SessionStoreSecure:   getEnvBool("SESSION_STORE_SECURE"),

		// Oauth
		// Google
		GoogleClientID:     getEnvString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnvString("GOOGLE_CLIENT_SECRET"),
		// Discord
		DiscordClientID:     getEnvString("DISCORD_CLIENT_ID"),
		DiscordClientSecret: getEnvString("DISCORD_CLIENT_SECRET"),

		// Stripe
		StripePublicKey: getEnvString("STRIPE_PUBLIC_KEY"),
		StripeSecretKey: getEnvString("STRIPE_SECRET_KEY"),
	}
}

// We might wantto fatal error if env not found
func getEnvString(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Printf("Env key not found: %s. Setting to an empty string.\n", key)
	return ""
}

func getEnvBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Printf("Error parsing %s as bool: %s\n", key, err)
			return false
		}
		return boolValue
	}

	fmt.Printf("Env key not found: %s. Setting to false.\n", key)
	return false
}
