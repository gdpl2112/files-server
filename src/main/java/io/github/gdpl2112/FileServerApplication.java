package io.github.gdpl2112;

import io.github.kloping.file.FileUtils;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.ServletComponentScan;

import java.io.File;
import java.lang.management.ManagementFactory;

@SpringBootApplication
@ServletComponentScan
public class FileServerApplication {
    public static void main(String[] args) {
        String name = ManagementFactory.getRuntimeMXBean().getName();
        String pid = name.split("@")[0];
        FileUtils.putStringInFile(pid, new File("./fs.pid"));
        SpringApplication.run(FileServerApplication.class, args);
    }
}