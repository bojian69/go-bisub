package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Security  SecurityConfig  `mapstructure:"security"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Redis     RedisConfig     `mapstructure:"redis"`
	WebUI     WebUIConfig     `mapstructure:"web_ui"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
}

type ServerConfig struct {
	Port      int           `mapstructure:"port"`
	Timeout   time.Duration `mapstructure:"timeout"`
	RateLimit int           `mapstructure:"rate_limit"`
}

type DatabaseConfig struct {
	Primary     DBConfig            `mapstructure:"primary"`
	DataSources map[string]DBConfig `mapstructure:"data_sources"`
}

type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type SecurityConfig struct {
	JWTSecret       string   `mapstructure:"jwt_secret"`
	AllowedSQLTypes []string `mapstructure:"allowed_sql_types"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type WebUIConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type SnowflakeConfig struct {
	NodeID int64 `mapstructure:"node_id"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
