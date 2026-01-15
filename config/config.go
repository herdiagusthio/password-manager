package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string `mapstructure:"SERVER_PORT"`
	DBSource         string `mapstructure:"DB_SOURCE"`
	RedisAddr        string `mapstructure:"REDIS_ADDR"`
	EncryptionKey    string `mapstructure:"ENCRYPTION_KEY"` // Master key for backup encryption
	GoogleClientID   string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	SessionSecret      string `mapstructure:"SESSION_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	viper.SetDefault("SERVER_PORT", ":8080")
	viper.SetDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/callback")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
			return
		}
		// Logic to allow running without config file if env vars are set
		log.Println("No config file found, relying on environment variables")
	}

	err = viper.Unmarshal(&config)
	return
}
