package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gdpl2112/files-server/model"
)

// AuthService 认证服务
type AuthService struct {
	AuthServerUrl string
	AppID         string
	AppSecret     string
	userCache     map[string]*model.User
	mutex         sync.RWMutex
}

// NewAuthService 创建认证服务实例
func NewAuthService(authServerUrl, appID, appSecret string) *AuthService {
	service := &AuthService{
		AuthServerUrl: authServerUrl,
		AppID:         appID,
		AppSecret:     appSecret,
		userCache:     make(map[string]*model.User),
	}
	service.LoadUser()
	return service
}

// LoadUser 加载用户缓存
func (s *AuthService) LoadUser() {
	filePath := "./data/user.json"
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(filePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("创建用户缓存目录失败: %v\n", err)
			}
			_, err = os.Create(filePath)
			if err != nil {
				fmt.Printf("创建用户缓存文件失败: %v\n", err)
			}
		}
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("读取用户缓存文件失败: %v\n", err)
		return
	}

	if len(bytes) == 0 {
		return
	}

	var users []*model.User
	if err := json.Unmarshal(bytes, &users); err != nil {
		fmt.Printf("解析用户缓存文件失败: %v\n", err)
		return
	}

	s.mutex.Lock()
	for _, user := range users {
		s.userCache[user.AccessToken] = user
	}
	s.mutex.Unlock()
}

// SaveCache 保存用户缓存
func (s *AuthService) SaveCache() {
	s.mutex.RLock()
	users := make([]*model.User, 0, len(s.userCache))
	for _, user := range s.userCache {
		users = append(users, user)
	}
	s.mutex.RUnlock()

	bytes, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("序列化用户缓存失败: %v\n", err)
		return
	}

	filePath := "./data/user.json"
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("创建用户缓存目录失败: %v\n", err)
		return
	}

	if err := ioutil.WriteFile(filePath, bytes, 0644); err != nil {
		fmt.Printf("写入用户缓存文件失败: %v\n", err)
		return
	}
}

// HandleAuthCallback 处理授权回调
func (s *AuthService) HandleAuthCallback(code string) *model.User {
	user := s.GetUserInfo(code)
	if user != nil {
		user.StorageLimit = 500 * 1024 * 1024 // 500MB
		user.UsedStorage = 0
		user.LoginTime = time.Now()

		s.mutex.Lock()
		s.userCache[user.AccessToken] = user
		s.mutex.Unlock()

		s.SaveCache()
	}
	return user
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(accessToken string) *model.User {
	url := fmt.Sprintf("%s/auth/app/user?app_secret=%s&user_code=%s", s.AuthServerUrl, s.AppSecret, accessToken)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求授权服务器失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("授权服务器返回错误状态码: %d\n", resp.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体失败: %v\n", err)
		return nil
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("解析响应体失败: %v\n", err)
		return nil
	}

	userID, ok := response["user_id"].(string)
	if !ok {
		fmt.Println("响应体中没有user_id字段")
		return nil
	}

	username := "Unknown"
	if nickname, ok := response["nickname"].(string); ok {
		username = nickname
	}

	return &model.User{
		UserID:       userID,
		Username:     username,
		AccessToken:  accessToken,
		LoginTime:    time.Now(),
		StorageLimit: 500 * 1024 * 1024,
		UsedStorage:  0,
	}
}

// UpdateUserStorage 更新用户存储使用情况
func (s *AuthService) UpdateUserStorage(userID string, size int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, user := range s.userCache {
		if user.UserID == userID {
			user.UsedStorage += size
			s.SaveCache()
			break
		}
	}
}

// GetUserByAccessToken 根据访问令牌获取用户
func (s *AuthService) GetUserByAccessToken(accessToken string) *model.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.userCache[accessToken]
}
