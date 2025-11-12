package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ip-query-system/internal/api"
	"ip-query-system/internal/config"
	"ip-query-system/internal/database"
	"ip-query-system/internal/ipdb"
	"ip-query-system/internal/middleware"
)

func main() {
	configPath := flag.String("config", "./config.yml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化MySQL数据库（同时用于IP查询）
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("初始化MySQL失败: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("关闭MySQL失败: %v", err)
		} else {
			log.Println("MySQL已关闭")
		}
	}()

	// 初始化IP查询器（基于MySQL）
	ipDB, err := ipdb.Init(db, cfg)
	if err != nil {
		log.Fatalf("初始化IP查询器失败: %v", err)
	}
	defer func() {
		if err := ipDB.Close(); err != nil {
			log.Printf("关闭IP查询器失败: %v", err)
		}
	}()

	// 初始化路由
	r := api.NewRouter(cfg, db, ipDB)

	// 设置中间件
	middleware.Setup(r, cfg)

	// 配置HTTP服务器
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动成功，监听端口: %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到关闭信号，开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Println("服务器已正常关闭")
}