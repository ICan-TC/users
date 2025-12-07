package config

import (
	"fmt"
	"os"

	cfg "github.com/ICan-TC/lib/config"
	"github.com/ICan-TC/lib/logging"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// --- Domain Sub-Structs ---
type ServerConfig struct {
	Port      int    `flag:"port" env:"SERVICE_PORT" yaml:"port" default:"8888" validate:"min=1,max=65535"`
	LogLevel  string `flag:"log_level" env:"LOG_LEVEL" yaml:"log_level" default:"info" validate:"oneof=debug info warn error"`
	LogFormat string `flag:"log_format" env:"LOG_FORMAT" yaml:"log_format" default:"text" validate:"oneof=text json"`
}

type DBConfig struct {
	Host     string `flag:"db_host" env:"DB_HOST" yaml:"db_host" validate:"required"`
	Port     int    `flag:"db_port" env:"DB_PORT" yaml:"db_port" validate:"min=1,max=65535"`
	Username string `flag:"db_username" env:"DB_USERNAME" yaml:"db_username" validate:"required"`
	Password string `flag:"db_password" env:"DB_PASSWORD" yaml:"db_password" validate:"required"`
	Name     string `flag:"db_name" env:"DB_NAME" yaml:"db_name" validate:"required"`
	SSL      bool   `flag:"db_ssl" env:"DB_SSL" yaml:"db_ssl"`
}

type AuthConfig struct {
	Secret          string `flag:"auth_secret" env:"AUTH_SECRET" yaml:"auth_secret" validate:"required"`
	AccessTokenTTL  int    `flag:"auth_access_token_ttl" env:"AUTH_ACCESS_TOKEN_TTL" yaml:"auth_access_token_ttl" validate:"min=1,max=86400"`
	RefreshTokenTTL int    `flag:"auth_refresh_token_ttl" env:"AUTH_REFRESH_TOKEN_TTL" yaml:"auth_refresh_token_ttl" validate:"min=1,max=86400"`
	RateLimit       int    `flag:"auth_rate_limit" env:"AUTH_RATE_LIMIT" yaml:"auth_rate_limit" validate:"min=1,max=1000"`
}

// --- Main Config Struct ---
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Auth   AuthConfig
}

var (
	config Config
	loaded bool
)

// --- Loader ---
func Load() {
	l := logging.L()
	l.Debug().Msg("loading config")
	if loaded {
		return
	}
	loaded = true

	v := viper.New()

	// Bind all config fields recursively
	cfg.BindConfigStruct(v, &config.Server, "server")
	cfg.BindConfigStruct(v, &config.DB, "db")
	cfg.BindConfigStruct(v, &config.Auth, "auth")

	// Bind CLI flags
	pflag.String("config", "", "Path to config file or directory")
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Load ENV variables
	v.AutomaticEnv()

	// Determine config file path from CLI or ENV, default to current directory
	configPath := v.GetString("config")
	if configPath == "" {
		configPath = v.GetString("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "."
	}
	l.Debug().Msgf("config path %s", configPath)

	if fi, err := os.Stat(configPath); err == nil && !fi.IsDir() {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(configPath)
	}
	if err := v.ReadInConfig(); err != nil {
		l.Warn().Msgf("failed to read config file: %v", err)
	}

	// Unmarshal to config struct
	if err := v.Unmarshal(&config); err != nil {
		l.Warn().Msgf("failed to unmarshal config: %v", err)
	}

	l.Debug().Any("server", config.Server).Msg("validate config struct")
	l.Debug().Any("db", config.DB).Msg("validate config struct")

	// Validate config
	l.Debug().Any("server", config.Server).Msg("validate config struct")
	if err := cfg.ValidateConfigStruct(&config.Server); err != nil {
		panic(fmt.Sprintf("Config for server validation error: %v", err))
	}
	if err := cfg.ValidateConfigStruct(&config.DB); err != nil {
		panic(fmt.Sprintf("Config for DB validation error: %v", err))
	}
	if err := cfg.ValidateConfigStruct(&config.Auth); err != nil {
		panic(fmt.Sprintf("config AuthConfig validation error: %v", err))
	}
}

// Get returns the global config
func Get() Config {
	return config
}
