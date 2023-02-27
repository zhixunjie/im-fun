我的实现

任务清单：
1. 注册、反注册
2. 心跳上报
3. 服务发现：HTTP调用、GRPC调用（google：GRPC的服务发现机制）

etcd：
微服务上报心跳（Lease续约） + 微服务Watch（HTTP Stream）

consul：
agent做健康检查（HTTP/TCP检查） + 微服务Watch（长轮询 Watch）