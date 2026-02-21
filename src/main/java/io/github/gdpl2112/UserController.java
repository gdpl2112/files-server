package io.github.gdpl2112;

import io.github.gdpl2112.model.StorageInfo;
import io.github.gdpl2112.model.User;
import io.github.gdpl2112.model.UserFile;
import io.github.gdpl2112.service.AuthService;
import io.github.gdpl2112.service.FileService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.Resource;
import org.springframework.core.io.UrlResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import jakarta.servlet.http.HttpSession;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private AuthService authService;

    @Autowired
    private FileService fileService;

    @Value("${file.upload-dir}")
    private String uploadDir;

    /**
     * 获取用户信息
     *
     * @param session HTTP会话
     * @return 用户信息
     */
    @GetMapping("/info")
    public ResponseEntity<User> getUserInfo(HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user != null) {
            return ResponseEntity.ok(user);
        } else {
            return ResponseEntity.status(401).build();
        }
    }

    /**
     * 获取存储使用情况
     *
     * @param session HTTP会话
     * @return 存储信息
     */
    @GetMapping("/storage")
    public ResponseEntity<StorageInfo> getStorageInfo(HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user == null) {
            return ResponseEntity.status(401).build();
        }

        StorageInfo storageInfo = new StorageInfo();
        long limit = user.getStorageLimit();
        long used = fileService.calculateFolderSize(user);
        long remaining = limit - used;
        double percentage = limit > 0 ? (double) used / limit * 100 : 0;

        storageInfo.setLimit(limit);
        storageInfo.setUsed(used);
        storageInfo.setRemaining(remaining);
        storageInfo.setPercentage(percentage);
        storageInfo.setLimitFormatted(StorageInfo.formatFileSize(limit));
        storageInfo.setUsedFormatted(StorageInfo.formatFileSize(used));
        storageInfo.setRemainingFormatted(StorageInfo.formatFileSize(remaining));

        return ResponseEntity.ok(storageInfo);
    }

    /**
     * 获取用户文件列表
     *
     * @param session HTTP会话
     * @return 文件列表
     */
    @GetMapping("/files")
    public ResponseEntity<List<UserFile>> getUserFiles(HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user == null) {
            return ResponseEntity.status(401).build();
        }

        List<UserFile> files = fileService.getUserFiles(user);
        return ResponseEntity.ok(files);
    }

    /**
     * 检查文件是否存在
     *
     * @param path 文件路径
     * @return 是否存在
     */
    @GetMapping("/exits")
    public ResponseEntity<Boolean> checkFileExists(@RequestParam String path,HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user == null) {
            return ResponseEntity.status(401).build();
        }
        File file = new File(user.getPathDir(uploadDir), path);
        return ResponseEntity.ok(file.exists());
    }

    /**
     * 下载文件
     *
     * @param filename 文件名
     * @return 文件资源
     */
    @GetMapping("/download/{filename:.+}")
    public ResponseEntity<Resource> downloadFile(HttpSession session, @PathVariable String filename) {
        User user = (User) session.getAttribute("user");
        if (user == null) return ResponseEntity.status(401).build();
        Path file = Paths.get(user.getPathDir(uploadDir)).resolve(filename);
        try {
            Resource resource = new UrlResource(file.toUri());
            if (resource.exists() || resource.isReadable()) {
                return ResponseEntity.ok()
                        .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + resource.getFilename() + "\"")
                        .body(resource);
            } else {
                return ResponseEntity.notFound().build();
            }
        } catch (IOException e) {
            return ResponseEntity.internalServerError().build();
        }
    }

    // 删除
    @DeleteMapping("/delete/{filename:.+}")
    public ResponseEntity<String> deleteFile(@PathVariable String filename, HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user == null) {
            return ResponseEntity.status(401).body("未登录");
        }

        File file = new File(user.getPathDir(uploadDir), filename);
        if (!file.exists()) {
            return ResponseEntity.notFound().build();
        }

        if (!file.delete()) {
            return ResponseEntity.status(500).body("删除文件失败");
        }

        return ResponseEntity.ok("删除成功");
    }

    /**
     * 上传文件
     *
     * @param file    上传的文件
     * @param session HTTP会话
     * @return 上传结果
     */
    @PostMapping("/upload")
    public ResponseEntity<String> uploadFile(@RequestParam("file") MultipartFile file, HttpSession session) {
        User user = (User) session.getAttribute("user");
        if (user == null) {
            return ResponseEntity.status(401).body("未登录");
        }

        if (file.isEmpty()) {
            return ResponseEntity.badRequest().body("文件为空");
        }

        try {
            // 检查文件大小是否超过用户存储限制
            long fileSize = file.getSize();
            long usedStorage = user.getUsedStorage();
            long storageLimit = user.getStorageLimit();

            if (usedStorage + fileSize > storageLimit) {
                return ResponseEntity.badRequest().body("存储空间不足");
            }

            // 保存文件到用户目录（这里简化处理，实际可能需要按用户ID分目录）
            File dir = new File(user.getPathDir(uploadDir));
            if (!dir.exists()) {
                dir.mkdirs();
            }
            String fileName = file.getOriginalFilename();
            if (fileName == null || fileName.isEmpty()) {
                fileName = System.currentTimeMillis() + ".dat";
            }

            File dest = new File(dir, fileName);
            byte[] bytes = file.getBytes();
            Files.write(dest.toPath(), bytes);

            // 更新用户已使用存储空间
            authService.updateUserStorage(user.getUserId(), fileSize);

            return ResponseEntity.ok("文件上传成功");
        } catch (Exception e) {
            return ResponseEntity.internalServerError().body("上传失败: " + e.getMessage());
        }
    }
}