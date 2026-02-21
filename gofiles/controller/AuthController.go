package controller

import (
	"fmt"
	"net/http"

	"github.com/gdpl2112/files-server/model"
	"github.com/gdpl2112/files-server/service"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	AuthService   *service.AuthService
	AuthServerUrl string
	AppID         string
	RedirectUri   string
}

// NewAuthController 创建认证控制器实例
func NewAuthController(authService *service.AuthService, authServerUrl, appID, redirectUri string) *AuthController {
	return &AuthController{
		AuthService:   authService,
		AuthServerUrl: authServerUrl,
		AppID:         appID,
		RedirectUri:   redirectUri,
	}
}

// Register 注册路由
func (c *AuthController) Register(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/login", c.Login)
		authGroup.GET("/callback", c.Callback)
		authGroup.GET("/user", c.GetCurrentUser)
		authGroup.POST("/logout", c.Logout)
	}
}

// Login 重定向到授权服务器登录
func (c *AuthController) Login(ctx *gin.Context) {
	authorizeUrl := fmt.Sprintf("%s/authc?app_id=%s&redirect_uri=%s", c.AuthServerUrl, c.AppID, c.RedirectUri)
	ctx.Redirect(http.StatusFound, authorizeUrl)
}

// Callback 处理授权服务器回调
func (c *AuthController) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "缺少授权码"})
		return
	}

	user := c.AuthService.HandleAuthCallback(code)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "登录失败"})
		return
	}

	// 将用户信息存储到会话中
	ctx.Set("user", user)
	ctx.Redirect(http.StatusFound, "/")
}

// GetCurrentUser 获取当前用户信息
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
		return
	}

	currentUser, ok := user.(*model.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "用户信息类型错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userId":       currentUser.UserID,
		"username":     currentUser.Username,
		"accessToken":  currentUser.AccessToken,
		"storageLimit": currentUser.StorageLimit,
		"usedSpace":    currentUser.UsedStorage,
	})
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 清除会话中的用户信息
	ctx.Set("user", nil)
	ctx.JSON(http.StatusOK, gin.H{"message": "已登出"})
}
