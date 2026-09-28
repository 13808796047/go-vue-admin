package main

import (
	"log"

	"go-vue-admin/config"
	"go-vue-admin/database"
	"go-vue-admin/router"
	"go-vue-admin/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 初始化数据库（自动迁移 + 内置 admin/123456）
	if err := database.Init(cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	svc := service.New(cfg, database.DB)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 8 << 20

	router.Register(r, service.Deps{Svc: svc, JWTSecret: cfg.JWTSecret}, database.DB)

	addr := cfg.Address()
	log.Printf("后端服务启动: %s  数据库: %s", addr, cfg.DBPath)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
