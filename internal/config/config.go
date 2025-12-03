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
	Level          string `mapstructure:"level"`
	Format         string `mapstructure:"format"`
	FileLogEnabled bool   `mapstructure:"file_log_enabled"`
	FileLogDir     string `mapstructure:"file_log_dir"`
	LogRequestBody bool   `mapstructure:"log_request_body"`
	LogResponseBody bool  `mapstructure:"log_response_body"`
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

	// 启用环境变量支持
	viper.AutomaticEnv()
	viper.SetEnvPrefix("") // 不使用前缀
	
	// 绑定环境变量到配置项
	// 数据库配置
	viper.BindEnv("database.primary.host", "DB_HOST")
	viper.BindEnv("database.primary.port", "DB_PORT")
	viper.BindEnv("database.primary.username", "DB_USER")
	viper.BindEnv("database.primary.password", "DB_PASSWORD")
	viper.BindEnv("database.primary.database", "DB_NAME")
	
	// 数据源配置
	viper.BindEnv("database.data_sources.default.host", "DATASOURCE_DEFAULT_HOST")
	viper.BindEnv("database.data_sources.default.port", "DATASOURCE_DEFAULT_PORT")
	viper.BindEnv("database.data_sources.default.username", "DATASOURCE_DEFAULT_USER")
	viper.BindEnv("database.data_sources.default.password", "DATASOURCE_DEFAULT_PASS")
	viper.BindEnv("database.data_sources.default.database", "DATASOURCE_DEFAULT_NAME")
	
	viper.BindEnv("database.data_sources.dbcfg_adb_uhomes.host", "DATASOURCE_UHOMES_HOST")
	viper.BindEnv("database.data_sources.dbcfg_adb_uhomes.port", "DATASOURCE_UHOMES_PORT")
	viper.BindEnv("database.data_sources.dbcfg_adb_uhomes.username", "DATASOURCE_UHOMES_USER")
	viper.BindEnv("database.data_sources.dbcfg_adb_uhomes.password", "DATASOURCE_UHOMES_PASS")
	viper.BindEnv("database.data_sources.dbcfg_adb_uhomes.database", "DATASOURCE_UHOMES_NAME")
	
	// Redis 配置
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	
	// JWT 配置
	viper.BindEnv("security.jwt_secret", "JWT_SECRET")
	
	// 日志配置
	viper.BindEnv("logging.level", "LOG_LEVEL")
	
	// Snowflake 配置
	viper.BindEnv("snowflake.node_id", "SNOWFLAKE_NODE_ID")
	
	// 服务器配置
	viper.BindEnv("server.port", "SERVER_PORT")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
