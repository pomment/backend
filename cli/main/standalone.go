package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pomment/pomment"
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/dao"
	"log"
	"path"
)

func StartStandaloneServer(basePath string) error {
	// 读取配置文件
	log.Println("Initializing configuration...")
	dao.InitDataBasePath(path.Join(basePath))
	err := config.InitConfig(basePath)

	if err != nil {
		log.Fatal(err)
	}

	// 初始化 HTTP 服务
	log.Println("Initializing http service...")
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// 设置 CORS
	if config.Content.System.DevelopCORSPolicy {
		log.Println("Initializing CORS policy...")
		engine.Use(cors.New(CorsConfig()))
	}

	// 初始化路由
	log.Println("Initializing routes...")
	pomment.InitStandaloneRoutes(engine, "/")

	// 初始化缓存
	log.Println("Initializing cache service...")
	if config.Content.Redis.Enabled {
		dao.ConnectToRedisServer(true, config.Content.Redis.Addr, config.Content.Redis.Password, config.Content.Redis.Database)
	}

	// 启动服务端
	log.Println("Starting server...")
	err = engine.Run(fmt.Sprintf("%s:%d", config.Content.System.Host, config.Content.System.Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped.")
	return nil
}

func CorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
	corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	return corsConf
}
