package io.github.gdpl2112.service;

import com.alibaba.fastjson.JSONObject;
import io.github.gdpl2112.model.User;
import io.github.kloping.file.FileUtils;
import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

@Slf4j
@Service
public class AuthService {

    @Value("${server.auth-server}")
    private String authServerUrl;

    @Value("${auth.app-id}")
    private String appId;

    @Value("${auth.app-secret}")
    private String appSecret;

    /**
     * access_token -> user
     */
    private final Map<String, User> userCache = new HashMap<>();

    private void saveCache(Map<String, User> userCache) {
        FileUtils.putStringInFile(JSONObject.toJSONString(userCache.values()), new File("./data/user.json"));
    }

    @PostConstruct
    public void loadUser() {
        File file = new File("./data/user.json");
        if (!file.exists()) {
            try {
                file.getParentFile().mkdirs();
                file.createNewFile();
            } catch (IOException e) {
               log.error("创建用户缓存文件失败: {}", e.getMessage(), e);
            }
        }
        String json = FileUtils.getStringFromFile(file.getAbsolutePath());
        if (json == null || json.isEmpty()) {
            return;
        }
        for (User user : JSONObject.parseArray(json, User.class)) {
            userCache.put(user.getAccessToken(), user);
        }
    }

    /**
     * 处理授权回调，获取用户信息
     *
     * @param code 授权码
     * @return 用户对象
     */
    public User handleAuthCallback(String code) {
        try {
            // 使用访问令牌获取用户信息
            User user = getUserInfo(code);
            if (user != null) {
                user.setStorageLimit(500 * 1024 * 1024);
                user.setUsedStorage(0);
                userCache.put(user.getAccessToken(), user);
                saveCache(userCache);
            }
            return user;
        } catch (Exception e) {
            log.error("获取用户信息失败: {}", e.getMessage(), e);
            return null;
        }
    }

    /**
     * 获取用户信息
     *
     * @param accessToken 访问令牌
     * @return 用户对象
     */
    private User getUserInfo(String accessToken) {
        try {
            RestTemplate restTemplate = new RestTemplate();
            // 发送请求到授权服务器获取用户信息
            String url = authServerUrl + "/auth/app/user" +
                    "?app_secret=" + appSecret +
                    "&user_code=" + accessToken;
            ResponseEntity<Map> response = restTemplate.getForEntity(url, Map.class);

            if (response.getStatusCode() == HttpStatus.OK) {
                Map responseBody = response.getBody();
                if (responseBody != null && responseBody.containsKey("user_id")) {
                    User user = new User();
                    user.setUserId((String) responseBody.get("user_id"));
                    user.setUsername(responseBody.getOrDefault("nickname", "Unknown").toString());
                    user.setAccessToken(accessToken);
                    user.setLoginTime(LocalDateTime.now());
                    return user;
                }
            }
        } catch (Exception e) {
            log.error("获取用户信息失败: {}", e.getMessage(), e);
        }
        return null;
    }

    /**
     * 更新用户已使用的存储空间
     *
     * @param userId 用户ID
     * @param size   文件大小
     */
    public void updateUserStorage(String userId, long size) {
        User user = userCache.get(userId);
        if (user != null) {
            user.setUsedStorage(user.getUsedStorage() + size);
            saveCache(userCache);
        }
    }
}