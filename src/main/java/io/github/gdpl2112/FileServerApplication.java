package io.github.gdpl2112;

import io.github.kloping.file.FileUtils;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.ServletComponentScan;
import org.springframework.web.bind.annotation.CrossOrigin;

import java.io.File;
import java.lang.management.ManagementFactory;

@Slf4j
@SpringBootApplication
@ServletComponentScan
@CrossOrigin
public class FileServerApplication {
    public static void main(String[] args) {
        String name = ManagementFactory.getRuntimeMXBean().getName();
        String pid = name.split("@")[0];
        FileUtils.putStringInFile(pid, new File("./fs.pid"));
        SpringApplication.run(FileServerApplication.class, args);
        log.info("--------FileServerApplication started compile at 25/08/24---------");
    }
}