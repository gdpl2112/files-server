package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gdpl2112/files-server/model"
	"github.com/gdpl2112/files-server/service"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	AuthService *service.AuthService
	FileService *service.FileService
	UploadDir   string
}

// NewUserController 创建用户控制器实例
func NewUserController(authService *service.AuthService, fileService *service.FileService, uploadDir string) *UserController {
	return &UserController{
		AuthService: authService,
		FileService: fileService,
		UploadDir:   uploadDir,
	}
}

// Register 注册路由
func (c *UserController) Register(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/info", c.GetUserInfo)
		userGroup.GET("/storage", c.GetStorageInfo)
		userGroup.GET("/files", c.GetUserFiles)
		userGroup.GET("/exits", c.CheckFileExists)
		userGroup.GET("/download/:filename", c.DownloadFile)
		userGroup.DELETE("/delete/:filename", c.DeleteFile)
		userGroup.POST("/upload", c.UploadFile)
	}
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, currentUser)
}

// GetStorageInfo 获取存储使用情况
func (c *UserController) GetStorageInfo(ctx *gin.Context) {
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

	used := c.FileService.CalculateUserFolderSize(currentUser)
	limit := currentUser.StorageLimit
	remaining := limit - used
	percentage := 0.0
	if limit > 0 {
		percentage = float64(used) / float64(limit) * 100
	}

	storageInfo := &model.StorageInfo{
		Limit:              limit,
		Used:               used,
		Remaining:          remaining,
		Percentage:         percentage,
		LimitFormatted:     model.FormatFileSize(limit),
		UsedFormatted:      model.FormatFileSize(used),
		RemainingFormatted: model.FormatFileSize(remaining),
	}

	ctx.JSON(http.StatusOK, storageInfo)
}

// GetUserFiles 获取用户文件列表
func (c *UserController) GetUserFiles(ctx *gin.Context) {
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

	files := c.FileService.GetUserFiles(currentUser)
	ctx.JSON(http.StatusOK, files)
}

// CheckFileExists 检查文件是否存在
func (c *UserController) CheckFileExists(ctx *gin.Context) {
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

	path := ctx.Query("path")
	filePath := filepath.Join(currentUser.GetPathDir(c.UploadDir), path)

	exists = false
	if _, err := os.Stat(filePath); err == nil {
		exists = true
	}

	ctx.JSON(http.StatusOK, exists)
}

// DownloadFile 下载文件
func (c *UserController) DownloadFile(ctx *gin.Context) {
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

	filename := ctx.Param("filename")
	filePath := filepath.Join(currentUser.GetPathDir(c.UploadDir), filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "文件不存在"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.File(filePath)
}

// DeleteFile 删除文件
func (c *UserController) DeleteFile(ctx *gin.Context) {
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

	filename := ctx.Param("filename")
	filePath := filepath.Join(currentUser.GetPathDir(c.UploadDir), filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "文件不存在"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除文件失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UploadFile 上传文件
func (c *UserController) UploadFile(ctx *gin.Context) {
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

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "文件为空"})
		return
	}

	// 检查文件大小是否超过用户存储限制
	fileSize := file.Size
	usedStorage := currentUser.UsedStorage
	storageLimit := currentUser.StorageLimit

	if usedStorage+fileSize > storageLimit {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "存储空间不足"})
		return
	}

	// 保存文件到用户目录
	dir := currentUser.GetPathDir(c.UploadDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建目录失败: " + err.Error()})
		return
	}

	fileName := file.Filename
	if fileName == "" {
		fileName = "file.dat"
	}

	dest := filepath.Join(dir, fileName)
	if err := ctx.SaveUploadedFile(file, dest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "上传失败: " + err.Error()})
		return
	}

	// 更新用户已使用存储空间
	c.AuthService.UpdateUserStorage(currentUser.UserID, fileSize)

	ctx.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}
