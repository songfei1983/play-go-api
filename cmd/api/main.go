package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/songfei1983/play-go-api/internal/app"
	"github.com/songfei1983/play-go-api/internal/config"
	"github.com/songfei1983/play-go-api/internal/server"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 初始化应用依赖
	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	// 初始化并启动服务器
	srv := server.New(app)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
