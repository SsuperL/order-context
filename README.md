# order-context
使用go实现DDD demo
采用grpc通信框架，并使用grpc-gateway支持REST方式

## version
go 1.17

## 数据库迁移 及部署
采用migrate工具进行迁移。
容器化后，采用`init container` + `job`的方式，等待数据库迁移完成再启动pod。
```
kubectl apply -f deploy/configmap.yaml
kubectl apply -f deploy/deployment.yaml
kubectl apply -f deploy/migrate_job.yaml
```