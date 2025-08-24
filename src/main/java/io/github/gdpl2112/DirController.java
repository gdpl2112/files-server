package io.github.gdpl2112;

import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.File;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.*;

/**
 * @author github kloping
 * @date 2025/8/24-12:19
 */
@RequestMapping("/dir")
@RestController
public class DirController {
    @Value("${file.upload-dir}")
    private String uploadDir;

    private static final String PARENT_HTML = "<!doctype html>\n" +
            "<html lang=\"zh-CN\">\n" +
            "<head>\n" +
            "    <meta charset=\"UTF-8\">\n" +
            "    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
            "    <title>KLOPING FILE PORT</title>\n" +
            "    <style>\n" +
            "        body {\n" +
            "            font-family: Arial, sans-serif;\n" +
            "            margin: 20px;\n" +
            "            line-height: 1.6;\n" +
            "        }\n" +
            "        table {\n" +
            "            width: 100%;\n" +
            "            border-collapse: collapse;\n" +
            "            margin-bottom: 20px;\n" +
            "        }\n" +
            "        th, td {\n" +
            "            border: 1px solid #ddd;\n" +
            "            padding: 8px;\n" +
            "            text-align: left;\n" +
            "        }\n" +
            "        th {\n" +
            "            background-color: #f2f2f2;\n" +
            "            font-weight: bold;\n" +
            "        }\n" +
            "        tr:nth-child(even) {\n" +
            "            background-color: #f9f9f9;\n" +
            "        }\n" +
            "        tr:hover {\n" +
            "            background-color: #f1f1f1;\n" +
            "        }\n" +
            "        a.folder {\n" +
            "            color: #8d02ff;\n" +
            "        }\n" +
            "        a.file {\n" +
            "            color: #5a3000;\n" +
            "        }\n" +
            "    </style>\n" +
            "</head>\n" +
            "<body>\n" +
            "<h1>/</h1>\n" +
            "<hr>\n" +
            "<table>\n" +
            "    <thead>\n" +
            "    <tr>\n" +
            "        <th>文件名</th>\n" +
            "        <th>大小</th>\n" +
            "        <th>修改日期</th>\n" +
            "        <th>类型</th>\n" +
            "    </tr>\n" +
            "    </thead>\n" +
            "    <tbody>\n" +
            "    </tbody>\n" +
            "</table>\n" +
            "<hr>\n" +
            "</body>\n" +
            "</html>";

    @RequestMapping
    public Object dir(
            HttpServletRequest request, HttpServletResponse response,
            @RequestParam(value = "path", required = false, defaultValue = "/") String path) {
        // 规范化路径
        String normalizedPath = normalizePath(path);

        File f0 = new File(uploadDir);
        File file = new File(f0, normalizedPath);

        Document document = org.jsoup.Jsoup.parse(PARENT_HTML);
        Element element = document.getElementsByTag("h1").get(0);
        element.text(path);
        List<File> files = new LinkedList<>(Arrays.asList(Objects.requireNonNull(file.listFiles())));
        files.sort((f1,f2)->{
            // 优先文件夹
            if (f1.isDirectory() && !f2.isDirectory()) {
                return -1;
            } else if (!f1.isDirectory() && f2.isDirectory()) {
                return 1;
            } else {
                // 都是文件夹或都是文件，按名称升序
                return f1.getName().compareToIgnoreCase(f2.getName());
            }
        });
        if (!path.equals("/")) {
            Element el0 = new Element("tr");
            Element el1 = new Element("td");
            el1.addClass("folder");

            Element ahref = new Element("a");
            ahref.attr("href", "/dir?path=" + getUpPath(path));
            ahref.text("../");
            el1.appendChild(ahref);
            el0.appendChild(el1);

            document.getElementsByTag("tbody").get(0).appendChild(el0);
        }
        for (File elFile : files) {
            String fname = elFile.getPath();
            fname = fname.replaceAll("\\\\", "/");
            fname = fname.replace(uploadDir, "");
            Element el0 = new Element("tr");

            Element el1 = new Element("td");

            Element ahref = new Element("a");
            ahref.addClass(elFile.isDirectory() ? "folder" : "file");
            String href = elFile.isDirectory() ?
                    "/dir?path=" + fname : fname;
            href = href.replaceAll("//", "/");
            ahref.attr("href", href);
            if (elFile.isFile()) ahref.attr("target", "_blank");
            ahref.text(elFile.getName());
            el1.appendChild(ahref);

            el0.appendChild(el1);

            Element el2 = new Element("td");
            el2.text(getFileSize(elFile));
            el0.appendChild(el2);

            Element el3 = new Element("td");
            el3.text(getTimed(elFile.lastModified()));
            el0.appendChild(el3);

            Element el4 = new Element("td");
            el4.text(elFile.isDirectory() ? "文件夹" : getLastName(elFile.getName()));
            el0.appendChild(el4);
            document.getElementsByTag("tbody").get(0).appendChild(el0);
        }
        String sb = document.html();
        response.setContentType("text/html");
        response.setCharacterEncoding("UTF-8");
        response.setStatus(200);
        try {
            response.getWriter().write(sb);
            response.getWriter().flush();
            response.getWriter().close();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }

    private static final SimpleDateFormat SF_0 = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss");

    public static String getTimed(long l) {
        return SF_0.format(new Date(l));
    }

    public static String getUpPath(String path) {
        String upPath = path;
        if (upPath.equals("/")) {
            return "/";
        } else {
            int n = upPath.lastIndexOf("/");
            upPath = upPath.substring(0, n);
            return upPath;
        }
    }

    public static String getLastName(String name) {
        int end0 = name.lastIndexOf(".");
        if (end0 < 0) {
            return "未知文件类型";
        }
        return name.substring(end0 + 1) + "文件";
    }

    public static String getFileSize(File listFile) {
        if (listFile.isFile()) {
            return formatBytes(listFile.length());
        } else return "-";
    }

    // 字节单位转换工具（如 B → GB）
    public static String formatBytes(long bytes) {
        if (bytes < 1024) return bytes + " B";
        int exp = (int) (Math.log(bytes) / Math.log(1024));
        char unit = "KMGTPE".charAt(exp - 1);
        return String.format("%.2f %sB", bytes / Math.pow(1024, exp), unit);
    }

    /**
     * 规范化路径，防止路径遍历攻击
     */
    public static String normalizePath(String path) {
        if (path == null || path.isEmpty()) {
            return "/";
        }

        // 替换反斜杠为正斜杠
        path = path.replace('\\', '/');

        // 处理相对路径
        String[] parts = path.split("/");
        List<String> normalizedParts = getNormalizedParts(parts);

        // 重新组合路径
        StringBuilder normalizedPath = new StringBuilder();
        for (String part : normalizedParts) {
            normalizedPath.append("/").append(part);
        }

        return normalizedPath.length() > 0 ? normalizedPath.toString() : "/";
    }

    private static List<String> getNormalizedParts(String[] parts) {
        List<String> normalizedParts = new ArrayList<>();

        for (String part : parts) {
            if (part.equals("..")) {
                // 回退上级目录，但不能超出根目录
                if (!normalizedParts.isEmpty() && !normalizedParts.get(normalizedParts.size() - 1).equals("..")) {
                    normalizedParts.remove(normalizedParts.size() - 1);
                } else {
                    // 如果试图超出根目录，则忽略
                    continue;
                }
            } else if (!part.equals(".") && !part.isEmpty()) {
                normalizedParts.add(part);
            }
        }
        return normalizedParts;
    }
}
