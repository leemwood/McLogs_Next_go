package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Cache    CacheConfig    `mapstructure:"cache"`
	AI       AIConfig       `mapstructure:"ai"`
	Patterns string         `mapstructure:"patterns"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type MongoDBConfig struct {
	URL        string `mapstructure:"url"`
	DB         string `mapstructure:"db"`
	Collection string `mapstructure:"collection"`
	Enabled    bool   `mapstructure:"enabled"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Enabled  bool   `mapstructure:"enabled"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type StorageConfig struct {
	CurrentID   string           `mapstructure:"current_id"`
	TTL         int64            `mapstructure:"time_to_live"`
	MaxLength   int              `mapstructure:"max_length"`
	MaxLines    int              `mapstructure:"max_lines"`
	Filesystem  FilesystemConfig `mapstructure:"filesystem"`
	MongoDB     EnabledConfig    `mapstructure:"mongodb"`
	Postgres    EnabledConfig    `mapstructure:"postgres"`
	Redis       EnabledConfig    `mapstructure:"redis"`
}

type FilesystemConfig struct {
	Path    string `mapstructure:"path"`
	Enabled bool   `mapstructure:"enabled"`
}

type EnabledConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type CacheConfig struct {
	Driver  string `mapstructure:"driver"`
	Enabled bool   `mapstructure:"enabled"`
}

type AIConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
