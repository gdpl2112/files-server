package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// DirController 目录控制器
type DirController struct {
	UploadDir string
}

// NewDirController 创建目录控制器实例
func NewDirController(uploadDir string) *DirController {
	return &DirController{
		UploadDir: uploadDir,
	}
}

// Register 注册路由
func (c *DirController) Register(router *gin.Engine) {
	router.GET("/dir", c.Dir)
}

// Dir 目录浏览
func (c *DirController) Dir(ctx *gin.Context) {
	path := ctx.DefaultQuery("path", "/")
	normalizedPath := c.normalizePath(path)

	filePath := filepath.Join(c.UploadDir, normalizedPath)
	fileInfo, err := os.Stat(filePath)
	if err != nil || !fileInfo.IsDir() {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "目录不存在"})
		return
	}

	entries, err := os.ReadDir(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "读取目录失败"})
		return
	}

	// 排序：文件夹优先，然后按名称排序
	sort.Slice(entries, func(i, j int) bool {
		iIsDir := entries[i].IsDir()
		jIsDir := entries[j].IsDir()

		if iIsDir && !jIsDir {
			return true
		}
		if !iIsDir && jIsDir {
			return false
		}

		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	// 生成HTML响应
	html := c.generateDirHTML(path, entries, normalizedPath)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, html)
}

// normalizePath 规范化路径
func (c *DirController) normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	normalizedParts := make([]string, 0)

	for _, part := range parts {
		if part == ".." {
			if len(normalizedParts) > 0 && normalizedParts[len(normalizedParts)-1] != ".." {
				normalizedParts = normalizedParts[:len(normalizedParts)-1]
			}
		} else if part != "." && part != "" {
			normalizedParts = append(normalizedParts, part)
		}
	}

	return strings.Join(normalizedParts, "/")
}

// getUpPath 获取上级目录路径
func (c *DirController) getUpPath(path string) string {
	if path == "/" {
		return "/"
	}

	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return "/"
	}

	return path[:lastSlash]
}

// getFileSize 获取文件大小
func (c *DirController) getFileSize(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "-"
	}

	if fileInfo.IsDir() {
		return "-"
	}

	return c.formatBytes(fileInfo.Size())
}

// formatBytes 格式化字节大小
func (c *DirController) formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	size := float64(bytes) / 1024
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

// getFileExtension 获取文件扩展名
func (c *DirController) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "未知文件类型"
	}
	return strings.TrimPrefix(ext, ".") + "文件"
}

// generateDirHTML 生成目录浏览HTML
func (c *DirController) generateDirHTML(path string, entries []os.DirEntry, normalizedPath string) string {
	var html strings.Builder

	html.WriteString(`<!doctype html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>KLOPING FILE PORT</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            line-height: 1.6;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        tr:hover {
            background-color: #f1f1f1;
        }
        a.folder {
            color: #8d02ff;
        }
        a.file {
            color: #5a3000;
        }
    </style>
</head>
<body>
`)

	html.WriteString(fmt.Sprintf("<h1>%s</h1>\n", path))
	html.WriteString("<hr>\n")
	html.WriteString(`<table>
    <thead>
    <tr>
        <th>文件名</th>
        <th>大小</th>
        <th>修改日期</th>
        <th>类型</th>
    </tr>
    </thead>
    <tbody>
`)

	// 添加上级目录链接
	if path != "/" {
		upPath := c.getUpPath(path)
		html.WriteString(fmt.Sprintf(`    <tr>
        <td><a href="/dir?path=%s" class="folder">../</a></td>
        <td>-</td>
        <td>-</td>
        <td>文件夹</td>
    </tr>
`, upPath))
	}

	// 添加文件和文件夹列表
	for _, entry := range entries {
		entryPath := filepath.Join(normalizedPath, entry.Name())
		webPath := strings.ReplaceAll(entryPath, "\\", "/")

		var href, cssClass, fileType string
		if entry.IsDir() {
			href = fmt.Sprintf("/dir?path=%s", webPath)
			cssClass = "folder"
			fileType = "文件夹"
		} else {
			href = "/" + webPath
			cssClass = "file"
			fileType = c.getFileExtension(entry.Name())
		}

		fileInfo, _ := entry.Info()
		modTime := fileInfo.ModTime().Format("2006/01/02 15:04:05")
		fileSize := c.getFileSize(filepath.Join(c.UploadDir, entryPath))

		targetBlank := ""
		if !entry.IsDir() {
			targetBlank = ` target="_blank"`
		}
		html.WriteString(fmt.Sprintf(`    <tr>
        <td><a href="%s" class="%s"%s>%s</a></td>
        <td>%s</td>
        <td>%s</td>
        <td>%s</td>
    </tr>
`,
			href, cssClass, targetBlank,
			entry.Name(), fileSize, modTime, fileType))
	}

	html.WriteString(`    </tbody>
</table>
<hr>
</body>
</html>`)

	return html.String()
}
