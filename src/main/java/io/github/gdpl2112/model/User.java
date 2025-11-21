package io.github.gdpl2112.model;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class User {
    private String userId;
    private String username;
    private String accessToken;
    private LocalDateTime loginTime;
    private long storageLimit; // 存储限制（字节）
    private long usedStorage;  // 已使用存储（字节）

    public String getPathDir(String dir) {
        return dir + "/users/" + userId;
    }
}