package model

import "time"

// User 用户模型
type User struct {
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"accessToken"`
	LoginTime    time.Time `json:"loginTime"`
	StorageLimit int64     `json:"storageLimit"`
	UsedStorage  int64     `json:"usedStorage"`
}

// GetPathDir 获取用户目录路径
func (u *User) GetPathDir(dir string) string {
	return dir + "/users/" + u.UserID
}
