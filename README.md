# micro-gin
## run: 
    cd docker && docker-compose up
## 相关服务访问地址    
    konga : http://127.0.0.1:1337/register
    
    jaeger: http://127.0.0.1:16686/search
    
    consul: http://127.0.0.1:8500/
    
    kibana: http://127.0.0.1:5601/
    
    test-compose-yml 目录中的consul-compose.yml不要执行，测试

## 计划整合：
    grpc-gateway    https://github.com/grpc-ecosystem/grpc-gateway
    熔断 降级 限流

