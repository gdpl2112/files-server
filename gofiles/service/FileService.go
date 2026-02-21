package service

import (
	"os"
	"path/filepath"

	"github.com/gdpl2112/files-server/model"
)

// FileService 文件服务
type FileService struct {
	UploadDir string
}

// NewFileService 创建文件服务实例
func NewFileService(uploadDir string) *FileService {
	return &FileService{
		UploadDir: uploadDir,
	}
}

// GetUserFiles 获取用户文件列表
func (s *FileService) GetUserFiles(user *model.User) []*model.UserFile {
	files := make([]*model.UserFile, 0)
	if user == nil {
		return files
	}

	userFolderPath := user.GetPathDir(s.UploadDir)
	if info, err := os.Stat(userFolderPath); err == nil && info.IsDir() {
		s.collectFiles(userFolderPath, &files, "")
	}

	return files
}

// collectFiles 递归收集文件
func (s *FileService) collectFiles(folderPath string, files *[]*model.UserFile, basePath string) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(folderPath, entry.Name())
		if entry.IsDir() {
			newBasePath := basePath
			if newBasePath != "" {
				newBasePath += "/"
			}
			newBasePath += entry.Name()
			s.collectFiles(entryPath, files, newBasePath)
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				continue
			}

			userFile := &model.UserFile{
				Name: entry.Name(),
				Size: fileInfo.Size(),
			}

			if basePath != "" {
				userFile.Path = basePath + "/" + entry.Name()
			} else {
				userFile.Path = entry.Name()
			}

			*files = append(*files, userFile)
		}
	}
}

// CalculateFolderSize 计算文件夹大小
func (s *FileService) CalculateFolderSize(folderPath string) int64 {
	var size int64

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0
	}

	return size
}

// CalculateUserFolderSize 计算用户文件夹大小
func (s *FileService) CalculateUserFolderSize(user *model.User) int64 {
	if user == nil {
		return 0
	}

	userFolderPath := user.GetPathDir(s.UploadDir)
	if info, err := os.Stat(userFolderPath); err == nil && info.IsDir() {
		return s.CalculateFolderSize(userFolderPath)
	}

	return 0
}
