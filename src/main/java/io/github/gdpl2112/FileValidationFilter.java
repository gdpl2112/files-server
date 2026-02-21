package io.github.gdpl2112;

import io.github.gdpl2112.model.User;
import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

@Slf4j
public class FileValidationFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {

        HttpServletRequest httpRequest = (HttpServletRequest) request;
        HttpServletResponse httpResponse = (HttpServletResponse) response;
        String requestURI = httpRequest.getRequestURI();
        if (requestURI.startsWith("/users")) {
            // 提取用户ID
            String userId = requestURI.substring("/users/".length(), requestURI.indexOf("/", "/users/".length()));
            HttpSession session = httpRequest.getSession(true);
            if (session != null) {
                User user = (User) session.getAttribute("user");
                if (user != null && user.getUserId().equals(userId)) {
                    // 验证通过，继续处理请求
                    chain.doFilter(request, response);
                    return;
                }
            }
            // 验证失败，返回错误信息
            httpResponse.setStatus(HttpServletResponse.SC_FORBIDDEN);
            httpResponse.setContentType("text/plain; charset=UTF-8");
            httpResponse.getOutputStream().write("无权访问".getBytes(StandardCharsets.UTF_8));
            httpResponse.getOutputStream().flush();
            httpResponse.getOutputStream().close();
            return;
        }

        // 继续执行下一个过滤器
        chain.doFilter(request, response);
    }

    @Override
    public void destroy() {
        // 销毁方法
    }
}