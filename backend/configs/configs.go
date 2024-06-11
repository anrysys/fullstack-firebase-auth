// Test: No test
// Description: This file contains the configuration for the application. It is responsible for loading the environment variables from the .env file.
package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName        string `mapstructure:"APP_NAME"`
	Protocol       string `mapstructure:"PROTOCOL"`
	Host           string `mapstructure:"HOST"`
	Port           string `mapstructure:"PORT"`
	ClientOrigin   string `mapstructure:"CLIENT_ORIGIN"`
	ApiVersion     string `mapstructure:"API_VERSION"`
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	DBSslMode      string `mapstructure:"POSTGRES_SSL_MODE"`
	TimeZone       string `mapstructure:"TIMEZONE"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ClientPort     string `mapstructure:"CLIENT_PORT"`
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      string `mapstructure:"REDIS_PORT"`
	RedisDB        int    `mapstructure:"REDIS_DB"`

	RedisUri string `mapstructure:"REDIS_URL"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtpPort     int    `mapstructure:"SMTP_PORT"`
	SmtpUser     string `mapstructure:"SMTP_USER"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`

	// Firebase configuration for Firebase Admin SDK (for App Check)
	// FirebaseAppID             string `mapstructure:"NEXT_PUBLIC_FIREBASE_APP_ID"`
	// FirebaseAPIKey            string `mapstructure:"NEXT_PUBLIC_FIREBASE_API_KEY"`
	// FirebaseAuthDomain        string `mapstructure:"NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN"`
	FirebaseDatabaseURL   string `mapstructure:"NEXT_PUBLIC_FIREBASE_DATABASE_URL"`
	FirebaseProjectID     string `mapstructure:"NEXT_PUBLIC_FIREBASE_PROJECT_ID"`
	FirebaseStorageBucket string `mapstructure:"NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET"`
	// FirebaseMessagingSenderID string `mapstructure:"NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID"`
	// FirebaseMeasurementID     string `mapstructure:"NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID"`
}

func LoadConfig() (config Config, err error) {

	// Load environment variables
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Load .env file
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Load .env.local file
	viper.SetConfigName(".env.local")
	viper.AddConfigPath("..")
	// Get only those variables that contain the prefix "NEXT_PUBLIC_FIREBASE_"
	viper.SetEnvPrefix("NEXT_PUBLIC_FIREBASE_")
	// viper.AutomaticEnv()
	err = viper.MergeInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
