package io.github.gdpl2112;

import jakarta.servlet.*;
import lombok.extern.slf4j.Slf4j;
import org.apache.catalina.connector.RequestFacade;

import java.io.IOException;

@Slf4j
public class LogFilter implements Filter {
    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        RequestFacade requestFacade = (RequestFacade) request;
        String ip = requestFacade.getHeader("x-forwarded-for");
        if (ip == null) ip = request.getRemoteAddr();
        ip = "0:0:0:0:0:0:0:1".equals(ip) ? "127.0.0.1" : ip;
        log.info("{}[{}]({})", ip, requestFacade.getMethod(), requestFacade.getRequestURL());
//        log.debug(com.alibaba.fastjson2.JSON.toJSONString(request.getParameterMap()));
        chain.doFilter(request, response);
    }
}