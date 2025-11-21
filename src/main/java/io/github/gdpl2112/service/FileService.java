package io.github.gdpl2112.service;

import io.github.gdpl2112.model.User;
import io.github.gdpl2112.model.UserFile;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

@Service
public class FileService {

    @Value("${file.upload-dir}")
    private String uploadDir;

    /**
     * 获取用户文件列表
     *
     * @param user 用户对象
     * @return 文件列表
     */
    public List<UserFile> getUserFiles(User user) {
        List<UserFile> files = new ArrayList<>();
        if (user == null) {
            return files;
        }

        // 构建用户文件夹路径（可以根据实际需求调整）
        String userFolderPath = user.getPathDir(uploadDir); // 这里简化处理，实际可能需要根据用户ID创建独立文件夹
        File userFolder = new File(userFolderPath);

        if (userFolder.exists() && userFolder.isDirectory()) {
            collectFiles(userFolder, files, "");
        }

        return files;
    }

    /**
     * 递归收集文件
     *
     * @param folder   文件夹
     * @param files    文件列表
     * @param basePath 基础路径
     */
    private void collectFiles(File folder, List<UserFile> files, String basePath) {
        File[] fileList = folder.listFiles();
        if (fileList != null) {
            for (File file : fileList) {
                if (file.isFile()) {
                    UserFile userFile = new UserFile();
                    userFile.setName(file.getName());
                    userFile.setSize(file.length());
                    userFile.setPath(basePath.isEmpty() ? file.getName() : basePath + "/" + file.getName());
                    files.add(userFile);
                } else if (file.isDirectory()) {
                    // 递归处理子文件夹
                    String newBasePath = basePath.isEmpty() ? file.getName() : basePath + "/" + file.getName();
                    collectFiles(file, files, newBasePath);
                }
            }
        }
    }

    /**
     * 计算文件夹大小
     *
     * @param folder 文件夹
     * @return 文件夹大小（字节）
     */
    public long calculateFolderSize(File folder) {
        long size = 0;
        File[] files = folder.listFiles();
        if (files != null) {
            for (File file : files) {
                if (file.isFile()) {
                    size += file.length();
                } else if (file.isDirectory()) {
                    size += calculateFolderSize(file);
                }
            }
        }
        return size;
    }

    public long calculateFolderSize(User user) {
        // 构建用户文件夹路径（可以根据实际需求调整）
        String userFolderPath = user.getPathDir(uploadDir); // 这里简化处理，实际可能需要根据用户ID创建独立文件夹
        File userFolder = new File(userFolderPath);
        if (userFolder.exists() && userFolder.isDirectory()) {
            return calculateFolderSize(userFolder);
        }
        return 0;
    }

}