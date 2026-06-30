package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DBConfig struct {
	DSN      string `mapstructure:"dsn"`
	PoolSize int    `mapstructure:"poolsize"`
}

type JWTConfig struct {
	AccessExpiry  time.Duration `mapstructure:"accessexpiry"`
	RefreshExpiry time.Duration `mapstructure:"refreshexpiry"`
}

type Config struct {
	Server          ServerConfig `mapstructure:"server"`
	DB              DBConfig     `mapstructure:"db"`
	JWT             JWTConfig    `mapstructure:"jwt"`
	EditWindowDays  int          `mapstructure:"editwindowdays"`
	LogLevel        string       `mapstructure:"loglevel"`
	Addr            string
	ShutdownTimeout time.Duration
	Version         string
}

func Default() Config {
	return Config{
		Server:          ServerConfig{Port: 8080},
		DB:              DBConfig{DSN: "postgres://postgres:postgres@localhost:5432/habitflow?sslmode=disable", PoolSize: 10},
		JWT:             JWTConfig{AccessExpiry: 15 * time.Minute, RefreshExpiry: 168 * time.Hour},
		EditWindowDays:  7,
		LogLevel:        "info",
		Addr:            ":8080",
		ShutdownTimeout: 30 * time.Second,
		Version:         "dev",
	}
}

func Load() Config {
	cfg := Default()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.BindEnv("SERVER_PORT")
	v.BindEnv("DB_DSN")
	v.BindEnv("DB_POOL_SIZE")
	v.BindEnv("JWT_ACCESS_EXPIRY")
	v.BindEnv("JWT_REFRESH_EXPIRY")
	v.BindEnv("EDIT_WINDOW_DAYS")
	v.BindEnv("LOG_LEVEL")
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unmarshal config: %w", err))
	}

	if port := getFirstString(v, "SERVER_PORT", "server.port"); port != "" {
		cfg.Server.Port = parsePort(port)
	}
	if dsn := getFirstString(v, "DB_DSN", "db.dsn"); dsn != "" {
		cfg.DB.DSN = dsn
	}
	if poolSize := getFirstInt(v, "DB_POOL_SIZE", "db.poolsize"); poolSize > 0 {
		cfg.DB.PoolSize = poolSize
	}
	if accessExpiry := getFirstDuration(v, "JWT_ACCESS_EXPIRY", "jwt.accessexpiry"); accessExpiry > 0 {
		cfg.JWT.AccessExpiry = accessExpiry
	}
	if refreshExpiry := getFirstDuration(v, "JWT_REFRESH_EXPIRY", "jwt.refreshexpiry"); refreshExpiry > 0 {
		cfg.JWT.RefreshExpiry = refreshExpiry
	}
	if editWindowDays := getFirstInt(v, "EDIT_WINDOW_DAYS", "editwindowdays"); editWindowDays > 0 {
		cfg.EditWindowDays = editWindowDays
	}
	if logLevel := getFirstString(v, "LOG_LEVEL", "loglevel"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		cfg.LogLevel = envLevel
	}

	cfg.Addr = fmt.Sprintf(":%d", cfg.Server.Port)
	return cfg
}

func parsePort(port string) int {
	port = strings.TrimSpace(port)
	if port == "" {
		return 8080
	}
	if strings.HasPrefix(port, ":") {
		port = strings.TrimPrefix(port, ":")
	}
	if port == "" {
		return 8080
	}
	var parsed int
	_, _ = fmt.Sscanf(port, "%d", &parsed)
	if parsed <= 0 {
		return 8080
	}
	return parsed
}

func getFirstString(v *viper.Viper, keys ...string) string {
	for _, key := range keys {
		if value := v.GetString(key); value != "" {
			return value
		}
	}
	return ""
}

func getFirstInt(v *viper.Viper, keys ...string) int {
	for _, key := range keys {
		if value := v.GetInt(key); value != 0 {
			return value
		}
	}
	return 0
}

func getFirstDuration(v *viper.Viper, keys ...string) time.Duration {
	for _, key := range keys {
		if value := v.GetDuration(key); value > 0 {
			return value
		}
	}
	return 0
}
