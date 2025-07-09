package io.github.gdpl2112;

import io.github.kloping.date.DateUtils;
import io.github.kloping.file.FileUtils;
import io.github.kloping.judge.Judge;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.Resource;
import org.springframework.core.io.UrlResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.UUID;

@Slf4j
@RestController
public class FileController {
    @Value("${file.upload-dir}")
    private String uploadDir;

    @Value("${server.host}")
    private String host;

    @Value("${server.port}")
    private Integer port;

    @PostMapping("/upload")
    public ResponseEntity<String> uploadFile(
            @RequestParam(name = "file") MultipartFile multipartFile
            , @RequestParam(name = "path", required = false) String path
            , @RequestParam(name = "name", required = false) String name
            , @RequestParam(name = "suffix", required = false) String suffix
    ) {
        if (multipartFile.isEmpty()) return ResponseEntity.badRequest().body("文件为空");
        try {
            if (io.github.kloping.judge.Judge.isEmpty(path)) {
                path = DateUtils.getYear() + "_" + DateUtils.getMonth();
            }
            File dir = new File(uploadDir, path);
            if (io.github.kloping.judge.Judge.isEmpty(name)) {
                name = DateUtils.getDay() + "-" + UUID.randomUUID();
                if (Judge.isEmpty(suffix)){
                    String ifn = multipartFile.getOriginalFilename();
                    suffix = ifn.substring(ifn.lastIndexOf("."));
                }
            }
            if (suffix != null) name = name + suffix;
            File dist = new File(dir.getPath(), name);
            FileUtils.testFile(dist);
            FileUtils.writeBytesToFile(multipartFile.getBytes(), dist);
            String outname = dist.getPath().replaceAll("\\\\", "/");
            outname = outname.replace(uploadDir, "");
            return ResponseEntity.ok(host + ":" + port + outname);
        } catch (Exception e) {
            log.error(e.getMessage(), e);
            return ResponseEntity.internalServerError().body("上传失败: " + e.getMessage());
        }
    }

    @GetMapping("/download/{filename:.+}")
    public ResponseEntity<Resource> downloadFile(@PathVariable String filename) {
        Path file = Paths.get(uploadDir).resolve(filename);
        try {
            Resource resource = new UrlResource(file.toUri());
            return ResponseEntity.ok()
                    .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + resource.getFilename() + "\"")
                    .body(resource);
        } catch (IOException e) {
            return ResponseEntity.notFound().build();
        }
    }
}