package io.github.gdpl2112.model;

import lombok.Data;

@Data
public class UserFile {
    private String name;    // 文件名
    private long size;      // 文件大小（字节）
    private String path;    // 文件路径
}