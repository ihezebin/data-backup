# data-backup
🥚鸡蛋不能放在一个篮子里！

## 简介
数据备份！通过定义【数据源】和【存储目标】，创建【备份任务】关联数据信息，来实现数据定时备份和数据恢复。

> 目前已实现 MongoDB指定库表 和 Minio目录文件 定时备份到 OSS

# 能力
- 组件化抽象定义，方便扩展实现新的数据源和备份存储目标。

- 简易的配置：
```toml
[[mongo_sources]]
    id = "mongo:blog"
    dsn = "mongodb://root:root@127.0.0.1:27017/blog?authSource=admin"
    collections = ["article", "tag", "draft", "comment"]

[[minio_sources]]
    id = "minio:blog"
    dsn = "minio://xxx:xxx@127.0.0.1:9000/blog-minio"
    prefixes = ["preview", "content"]

[[oss_targets]]
    id = "cos:backup:mongo"
    dsn = "cos://xxx:xxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727"
    dir = "backup_mongo"
[[oss_targets]]
    id = "cos:backup:minio"
    dsn = "cos://xxx:xxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727"
    dir = "backup_minio"


[[tasks]]
    id = "664c72b790d71012f2753739"
    cron = "0 1 * * * *"
    source_id = "mongo:blog"
    target_id = "cos:backup:mongo"

[[tasks]]
    id = "664c72b790d71012f2753731"
    cron = "0 1 * * * *"
    source_id = "minio:blog"
    target_id = "cos:backup:minio"
```

- 备份失败告警通知，目前实现邮件告警
- 数据恢复
```curl
curl -X POST http://127.0.0.1:8080/restore/task \
     -H "Content-Type: application/json" \
     -d '{"task_id": "664c72b790d71012f2753739"}'
```

# 计划
- Web UI 操作页面
- 其他类型数据库数据备份，如 Mysql、Redis 等
- 探索非数据库类型的数据备份，如 Git 仓库项目代码等其他日常生活动需要预防数据意外丢失的场景