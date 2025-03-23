package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host string
		Port string
	}
	Server struct {
		Port string
	}
	Tracing struct {
		Endpoint string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// 从环境变量加载配置
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	cfg.Tracing.Endpoint = os.Getenv("TRACING_ENDPOINT")
	if cfg.Tracing.Endpoint == "" {
		cfg.Tracing.Endpoint = "jaeger:4317"
	}

	return cfg, nil
}

func (c *Config) GetDBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
