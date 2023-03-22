# Train-Tiktok

字节青训营第五期 大作业, 简易抖音后端设计

项目介绍地址：[不用最好的抖音极简版](https://bqz5oeip06.feishu.cn/docx/MYL6dzRkwoi7m2xivcZcjyspn0b)

use go-zero as microservice framework

## Ports

| port | service  | description      |
|:-----|:---------|:-----------------|
| 8888 | gateway  | Gateway          |
| 8081 | identity | identity service |
| 8082 | user     | user service     |
| 8083 | video    | video service    |
| 8084 | chat     | chat service     |

## Deploy

```bash
# 启动依赖 (etcd, redis, mysql)
docker-compose -f docker-compose.env.yaml up -d

# 启动服务
docker-compose -f docker-compose.yaml up -d

# 编译运行
docker-compose -f docker-compose.build.yaml up -d
```