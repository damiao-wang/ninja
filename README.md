# ninja
---
**ninja**是一个Go Web框架，目的是为了能够快速开发RestAPI，以下是ninja需要实现的功能。
>1. WebAPI采用AUTO的方式注册路由;
>2. 底层采用gorilla/mux和negroni;
>3. 实现log;
>4. 实现error追踪;
>5. 实现数据表的解析，以实现快速开发;
>6. 兼容GRPC，使WebAPI和GRPC采用同一个处理函数，使用RPC通信;

**task**
>2. 移植context
>3. 移植statsd
>4. 兼容下grpc （tommorrow）
>5. 移植ezcache
>6. 移植consul   负载均衡
>7. 学习部署
>8. 单元测试
