package config

import (
	"log"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config struct defines the application’s configuration schema
type Config struct {
	Database struct {
		Host     string `default:"127.0.0.1"`
		Port     string `default:"5432"`
		User     string
		Password string
		Name     string `default:"wallet_db"`
		SSLMode  string `default:"disable"`
	}

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

	if len(AppConfig.Concurrencies) == 0 {
		log.Fatal("no concurrency configurations found")
	}
}

// SetupDatabase initializes the database connection using GORM and the config values
func SetupDatabase() *gorm.DB {
	dsn := "host=" + AppConfig.Database.Host +
		" user=" + AppConfig.Database.User +
		" password=" + AppConfig.Database.Password +
		" dbname=" + AppConfig.Database.Name +
		" port=" + AppConfig.Database.Port +
		" sslmode=" + AppConfig.Database.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return db
}
