package model

// UserFile 用户文件模型
type UserFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}
