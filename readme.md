### 简易 File Server

> https://kloping.top/ 的文件服务器

接口:
- POST /upload
- GET  /dir

## 前端页面

项目包含以下前端页面：
- `index.html` - 主页，提供用户认证入口和文件操作界面
- `login.html` - 登录页面
- `user.html` - 用户中心，展示用户信息和文件管理功能

## 新增认证功能

现在支持用户认证和文件隔离功能：

### 认证接口:
- GET  /auth/login - 用户登录，重定向到授权服务器
- GET  /auth/callback - 授权回调处理
- GET  /auth/user - 获取当前用户信息
- POST /auth/logout - 用户登出

### 用户文件操作接口:
- POST /user/upload - 用户上传文件（需要登录）
- GET  /user/download/{filename} - 用户下载文件（需要登录）
- GET  /user/exits?path={文件路径} - 检查用户文件是否存在（需要登录）
- GET  /user/storage/info - 查看用户存储使用情况（需要登录）

### 功能特性:
1. 用户认证授权登录
2. 每个用户文件隔离，只能访问自己的文件
3. 文件上传和下载
4. 存储空间限制和监控（默认每个用户100MB存储空间）

### 配置说明:
在 `application.yml` 中配置以下参数：
```yaml
server:
  host: http://127.0.0.1      # 服务器主机地址
  port: 82                    # 服务器端口
  auth-server: http://127.0.0.1:81  # 授权服务器地址

auth:
  app-id: 101026453          # 应用ID
  app-secret: Ue7V4HmpMBXlQT2 # 应用密钥
  redirect-uri: ${server.host}:${server.port}/auth/callback # 回调地址
```

### 用户文件隔离:
每个用户的文件都存储在独立的目录中：
```
/files/user_{userId}/
```