package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileController 文件控制器
type FileController struct {
	UploadDir string
	Host      string
	Port      int
}

// NewFileController 创建文件控制器实例
func NewFileController(uploadDir, host string, port int) *FileController {
	return &FileController{
		UploadDir: uploadDir,
		Host:      host,
		Port:      port,
	}
}

// Register 注册路由
func (c *FileController) Register(router *gin.Engine) {
	router.POST("/upload", c.UploadFile)
	router.GET("/download/:filename", c.DownloadFile)
	router.GET("/exits", c.Exits)
	router.GET("/ping", c.Ping)
}

// UploadFile 上传文件
func (c *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "文件为空"})
		return
	}

	path := ctx.PostForm("path")
	name := ctx.PostForm("name")
	suffix := ctx.PostForm("suffix")

	if path == "" {
		path = fmt.Sprintf("%d_%d", time.Now().Year(), time.Now().Month())
	}

	dir := filepath.Join(c.UploadDir, path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建目录失败: " + err.Error()})
		return
	}

	if name == "" {
		name = fmt.Sprintf("%d-%s", time.Now().Day(), uuid.New().String())
		if suffix == "" {
			if file.Filename != "" {
				suffix = filepath.Ext(file.Filename)
			}
		}
	}

	if suffix != "" {
		name += suffix
	}

	dest := filepath.Join(dir, name)
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "打开文件失败: " + err.Error()})
		return
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建文件失败: " + err.Error()})
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "写入文件失败: " + err.Error()})
		return
	}

	outname := dest
	outname = outname[len(c.UploadDir):]
	outname = filepath.ToSlash(outname)

	ctx.JSON(http.StatusOK, fmt.Sprintf("%s:%d%s", c.Host, c.Port, outname))
}

// DownloadFile 下载文件
func (c *FileController) DownloadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filePath := filepath.Join(c.UploadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")
	ctx.File(filePath)
}

// Exits 检查文件是否存在
func (c *FileController) Exits(ctx *gin.Context) {
	path := ctx.Query("path")
	filePath := filepath.Join(c.UploadDir, path)

	exists := false
	if _, err := os.Stat(filePath); err == nil {
		exists = true
	}

	ctx.JSON(http.StatusOK, exists)
}

// Ping 心跳检测
func (c *FileController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
