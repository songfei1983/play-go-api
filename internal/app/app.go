package app

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/songfei1983/play-go-api/internal/config"
	"github.com/songfei1983/play-go-api/internal/handler"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB      *gorm.DB
	Redis   *redis.Client
	cleanup func()
}

func New(cfg *config.Config) (*App, error) {
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, err
	}

	// 初始化OpenTelemetry追踪器
	cleanup, err := initTracer(cfg.Tracing.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	return &App{
		DB:      db,
		Redis:   redisClient,
		cleanup: cleanup,
	}, nil
}

func (a *App) Close() {
	if sqlDB, err := a.DB.DB(); err == nil {
		sqlDB.Close() // nolint: errcheck
	}
	if a.Redis != nil {
		a.Redis.Close() // nolint: errcheck
	}
	if a.cleanup != nil {
		a.cleanup()
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.GetDBDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the User model
	if err := db.AutoMigrate(&handler.User{}); err != nil {
		return nil, err
	}

	return db, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.GetRedisAddr(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	return client, err
}
