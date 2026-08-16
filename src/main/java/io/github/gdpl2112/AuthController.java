package io.github.gdpl2112;

import io.github.gdpl2112.model.User;
import io.github.gdpl2112.service.AuthService;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/auth")
public class AuthController {
    
    @Autowired
    private AuthService authService;
    
    @Value("${server.auth-server}")
    private String authServerUrl;
    
    @Value("${auth.app-id}")
    private String appId;
    
    @Value("${auth.redirect-uri}")
    private String redirectUri;
    
    /**
     * 重定向到授权服务器进行登录
     * @return 重定向URL
     */
    @GetMapping("/login")
    public ResponseEntity<String> login(HttpServletResponse response) throws IOException {
        String authorizeUrl = String.format("%s/authc?app_id=%s&redirect_uri=%s", authServerUrl, appId, redirectUri);
        response.sendRedirect(authorizeUrl);
        return ResponseEntity.status(302).body(authorizeUrl);
    }
    
    /**
     * 处理授权服务器的回调
     * @param code 授权码
     * @param session HTTP会话
     * @return 登录结果
     */
    @GetMapping("/callback")
    public ResponseEntity<String> callback(@RequestParam("code") String code, HttpSession session, HttpServletResponse response) throws IOException {
        User user = authService.handleAuthCallback(code);
        if (user != null) {
            session.setAttribute("user", user);
            response.sendRedirect("/");
            return ResponseEntity.ok("登录成功，欢迎 " + user.getUsername());
        } else {
            return ResponseEntity.status(401).body("登录失败");
        }
    }
    
    /**
     * 获取当前用户信息
     * @param session HTTP会话
     * @return 当前用户信息
     */
    @GetMapping("/user")
    public ResponseEntity<Map<String, Object>> getCurrentUser(HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user != null) {
            Map<String, Object> userInfo = new HashMap<>();
            userInfo.put("userId", user.getUserId());
            userInfo.put("username", user.getUsername());
            userInfo.put("accessToken", user.getAccessToken());
            userInfo.put("storageLimit", user.getStorageLimit());
            userInfo.put("usedSpace", user.getUsedStorage());
            return ResponseEntity.ok(userInfo);
        } else {
            return ResponseEntity.status(401).build();
        }
    }
    
    /**
     * 用户登出
     * @param session HTTP会话
     * @return 登出结果
     */
    @PostMapping("/logout")
    public ResponseEntity<String> logout(HttpSession session) {
        session.invalidate();
        return ResponseEntity.ok("已登出");
    }
}