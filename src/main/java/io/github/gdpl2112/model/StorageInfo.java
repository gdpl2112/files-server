package io.github.gdpl2112.model;

import lombok.Data;

import java.text.DecimalFormat;

@Data
public class StorageInfo {
    private long limit;           // 存储限制（字节）
    private long used;            // 已使用存储（字节）
    private long remaining;        // 剩余存储（字节）
    private double percentage;     // 使用百分比
    private String limitFormatted;    // 格式化的存储限制
    private String usedFormatted;     // 格式化的已使用存储
    private String remainingFormatted; // 格式化的剩余存储


    /**
     * 格式化文件大小
     * @param bytes 字节数
     * @return 格式化后的文件大小字符串
     */
    public static String formatFileSize(long bytes) {
        if (bytes <= 0) return "0 Bytes";

        final String[] units = new String[]{"Bytes", "KB", "MB", "GB", "TB"};
        int digitGroups = (int) (Math.log10(bytes) / Math.log10(1024));
        return new DecimalFormat("#,##0.#").format(bytes / Math.pow(1024, digitGroups)) + " " + units[digitGroups];
    }
}