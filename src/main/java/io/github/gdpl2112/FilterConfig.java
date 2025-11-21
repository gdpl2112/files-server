package io.github.gdpl2112;

import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class FilterConfig {
    @Bean
    public FilterRegistrationBean<FileValidationFilter> fileValidationFilter() {
        FilterRegistrationBean<FileValidationFilter> bean = new FilterRegistrationBean<>();
        bean.setFilter(new FileValidationFilter());
        bean.addUrlPatterns("/*"); // 只拦截下载请求
        bean.setOrder(1); // 设置优先级高于日志过滤器
        return bean;
    }

    @Bean
    public FilterRegistrationBean<LogFilter> logFilter() {
        FilterRegistrationBean<LogFilter> bean = new FilterRegistrationBean<>();
        bean.setFilter(new LogFilter());
        bean.addUrlPatterns("/*"); // 拦截所有请求
        bean.setOrder(2); // 数字越小优先级越高
        return bean;
    }
}