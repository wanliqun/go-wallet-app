package config

import (
	"log"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Config struct defines the application’s configuration schema
type Config struct {
	Database DatabaseConfig

	Server struct {
		Port string `default:"8080"`
	}

	Concurrencies map[string]ConcurrencyConfig
}

type ConcurrencyConfig struct {
	Name      string
	Precision int
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig initializes Viper and loads the configuration from file and environment variables
func LoadConfig() {
	defaults.SetDefaults(&AppConfig)

	// Set the config file and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")

	// Read configuration from config.yaml if available
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Could not read config file: %v\n", err)
	}

	// Set up environment variable bindings (prefix with APP_)
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal the config into AppConfig struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

// SetupDatabase initializes the database connection using GORM and the config values
func SetupDatabase() *gorm.DB {
	return AppConfig.Database.MustOpenOrCreate()
}
