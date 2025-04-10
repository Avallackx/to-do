package config

import (
	"fmt"
	"strings"
	"time"

	"todo-app/internal/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warningf("%v", err)
	}
}

func Env() string {
	return viper.GetString("env")
}

func ServerPort() string {
	return viper.GetString("server_port")
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// PostgresHost :nodoc:
func PostgresHost() string {
	return viper.GetString("postgres.host")
}

// PostgresDatabase :nodoc:
func PostgresDatabase() string {
	return viper.GetString("postgres.database")
}

// PostgresUsername :nodoc:
func PostgresUsername() string {
	return viper.GetString("postgres.username")
}

// PostgresPassword :nodoc:
func PostgresPassword() string {
	return viper.GetString("postgres.password")
}

// PostgresSSLMode :nodoc:
func PostgresSSLMode() string {
	if viper.IsSet("postgres.sslmode") {
		return viper.GetString("postgres.sslmode")
	}
	return "disable"
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s@%s/%s?sslmode=%s",
		PostgresUsername(),
		PostgresHost(),
		PostgresDatabase(),
		PostgresSSLMode())
	// PostgresPassword(),

	// "postgres://postgres@db:5432/to-do-db?sslmode=disable"
}

// PostgresMaxIdleConns :nodoc:
func PostgresMaxIdleConns() int {
	if viper.GetInt("postgres.max_idle_conns") <= 0 {
		return DefaultPostgresMaxIdleConns
	}

	return viper.GetInt("postgres.max_idle_conns")
}

// PostgresMaxOpenConns :nodoc:
func PostgresMaxOpenConns() int {
	if viper.GetInt("postgres.max_open_conns") <= 0 {
		return DefaultPostgresMaxOpenConns
	}

	return viper.GetInt("postgres.max_open_conns")
}

// PostgresConnMaxLifetime :nodoc:
func PostgresConnMaxLifetime() time.Duration {
	cfg := viper.GetString("postgres.conn_max_lifetime")
	return utils.ParseDuration(cfg, DefaultPostgresConnMaxLifetime)
}

// PostgresPingInterval :nodoc:
func PostgresPingInterval() time.Duration {
	cfg := viper.GetString("postgres.ping_interval")
	return utils.ParseDuration(cfg, DefaultPostgresPingInterval)
}

// PostgresRetryAttempts :nodoc:
func PostgresRetryAttempts() int {
	if viper.GetInt("postgres.retry_attempts") > 0 {
		return viper.GetInt("postgres.retry_attempts")
	}

	return DefaultPostgresRetryAttempts
}

// RedisHost :nodoc:
func RedisHost() string {
	return viper.GetString("redis.host")
}

// RedisPassword :nodoc:
func RedisPassword() string {
	return viper.GetString("redis.password")
}

// RedisDB :nodoc:
func RedisDB() int {
	return viper.GetInt("redis.db")
}
