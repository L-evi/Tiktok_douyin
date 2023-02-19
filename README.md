# Train-Tiktok

字节青训营第五期 大作业, 简易抖音后端设计

use go-zero as microservice framework

## ports

| port | service  | description      |
|:-----|:---------|:-----------------|
| 8888 | gateway  | Gateway          |
| 8081 | identity | identity service |
| 8082 | user     | user service     |
| 8083 | video    | video service    |

## run

```bash
# 启动依赖 (etcd, redis, mysql)
docker-compose -f docker-compose.env.yaml up -d
# 启动服务
docker-compose -f docker-compose.yaml up -d
```