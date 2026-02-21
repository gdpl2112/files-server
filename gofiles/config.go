package main

import (
	"flag"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Host string
	Port int

	// 文件配置
	UploadDir string

	// 认证配置
	AuthServerUrl string
	AppID         string
	AppSecret     string
	RedirectUri   string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 从环境变量或命令行参数加载配置
	host := getEnv("SERVER_HOST", "localhost")
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	uploadDir := getEnv("FILE_UPLOAD_DIR", "./files")
	authServerUrl := getEnv("SERVER_AUTH_SERVER", "https://kloping.top")
	appID := getEnv("AUTH_APP_ID", "101610632")
	appSecret := getEnv("AUTH_APP_SECRET", "6mpmwW1axbbxyIT")
	redirectUri := getEnv("AUTH_REDIRECT_URI", "https://file.kloping.top/auth/callback")

	// 解析命令行参数
	flag.StringVar(&host, "host", host, "服务器主机名")
	flag.IntVar(&port, "port", port, "服务器端口")
	flag.StringVar(&uploadDir, "upload-dir", uploadDir, "文件上传目录")
	flag.StringVar(&authServerUrl, "auth-server", authServerUrl, "认证服务器地址")
	flag.StringVar(&appID, "app-id", appID, "应用ID")
	flag.StringVar(&appSecret, "app-secret", appSecret, "应用密钥")
	flag.StringVar(&redirectUri, "redirect-uri", redirectUri, "重定向URI")
	flag.Parse()

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic("创建上传目录失败: " + err.Error())
	}

	return &Config{
		Host:          host,
		Port:          port,
		UploadDir:     uploadDir,
		AuthServerUrl: authServerUrl,
		AppID:         appID,
		AppSecret:     appSecret,
		RedirectUri:   redirectUri,
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
