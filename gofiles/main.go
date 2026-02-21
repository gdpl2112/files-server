package main

import (
	"fmt"
	"os"

	"github.com/gdpl2112/files-server/controller"
	"github.com/gdpl2112/files-server/middleware"
	"github.com/gdpl2112/files-server/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config := LoadConfig()

	// 初始化服务
	authService := service.NewAuthService(config.AuthServerUrl, config.AppID, config.AppSecret)
	fileService := service.NewFileService(config.UploadDir)

	// 初始化控制器
	fileController := controller.NewFileController(config.UploadDir, config.Host, config.Port)
	dirController := controller.NewDirController(config.UploadDir)
	authController := controller.NewAuthController(authService, config.AuthServerUrl, config.AppID, config.RedirectUri)
	userController := controller.NewUserController(authService, fileService, config.UploadDir)

	// 创建Gin引擎
	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册中间件
	router.Use(middleware.LogFilter())
	router.Use(middleware.FileValidationFilter())

	// 注册路由
	fileController.Register(router)
	dirController.Register(router)
	authController.Register(router)
	userController.Register(router)

	// 静态文件服务
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	fmt.Printf("服务器启动在 %s\n", addr)
	if err := router.Run(addr); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}
