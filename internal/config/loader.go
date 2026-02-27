package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig загружает конфигурацию из YAML файла с подстановкой переменных окружения
func LoadConfig() AppConfig {
	cfg := AppConfig{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	// Enable environment variable substitution in config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		// Если файл не найден, используем значения по умолчанию
		return cfg
	}

	// Database configuration
	if viper.GetString("database.host") != "" && viper.GetString("database.port") != "" {
		connMaxLifetime := viper.GetDuration("database.conn_max_lifetime")
		if connMaxLifetime == 0 {
			connMaxLifetime = time.Hour
		}

		cfg.DB = &DBConfig{
			Driver: "mysql",
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				expandEnv(viper.GetString("database.user")),
				expandEnv(viper.GetString("database.password")),
				expandEnv(viper.GetString("database.host")),
				expandEnv(viper.GetString("database.port")),
				expandEnv(viper.GetString("database.name")),
			),
			MaxOpenConns: viper.GetInt("database.max_open_conns"),
			MaxIdleConns: viper.GetInt("database.max_idle_conns"),
			ConnMaxLife:  connMaxLifetime,
		}
	}

	// Redis configuration
	if viper.GetString("redis.addr") != "" {
		ttl := viper.GetDuration("cache.ttl")
		if ttl == 0 {
			ttl = 5 * time.Minute
		}

		cfg.Redis = &RedisConfig{
			Addr:     expandEnv(viper.GetString("redis.addr")),
			Password: expandEnv(viper.GetString("redis.password")),
			DB:       viper.GetInt("redis.db"),
			TTL:      ttl,
		}
	}

	// JWT configuration
	if viper.GetString("jwt.secret") != "" {
		expiration := viper.GetDuration("jwt.expiration")
		if expiration == 0 {
			expiration = 24 * time.Hour
		}

		cfg.JWT = &JWTConfig{
			Secret:     expandEnv(viper.GetString("jwt.secret")),
			Expiration: expiration,
		}
	}

	return cfg
}

func setDefaults() {
	// Database
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.port", "")
	viper.SetDefault("database.name", "")
	viper.SetDefault("database.user", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	// Redis
	viper.SetDefault("redis.addr", "")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT defaults
	viper.SetDefault("jwt.secret", "")
	viper.SetDefault("jwt.expiration", "24h")

	// Cache defaults
	viper.SetDefault("cache.ttl", "5m")
}

func expandEnv(s string) string {
	return os.ExpandEnv(s)
}
