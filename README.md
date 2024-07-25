# data-backup
🥚鸡蛋不能放在一个篮子里！

## 简介
数据备份！通过定义【数据源】和【备份存储目标】，创建【备份任务】关联数据信息，来实现数据定时备份和数据恢复。

> 目前已实现 MongoDB 指定库表定时备份到 OSS

# 能力
- 组件化抽象定义，方便扩展实现新的数据源和备份存储目标。

- 简易的配置：
```toml
[[mongo_sources]]
    key = "blog"
    dsn = "mongodb://root:root@mongo-sts-0.mongo.default:27017,mongo-sts-1.mongo.default:27017/blog?authSource=admin&replicaSet=hezebin"
    collections = ["article", "tag", "draft", "comment"]

[[oss_targets]]
    key = "cos"
    dsn = "cos://xxx:xxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727"
    dir = "backup"


[[tasks]]
    id = "664c72b790d71012f2753739"
    cron = "0 1 * * * *"
    source_type = "mongo"
    source_key = "blog"
    target_type = "oss"
    target_key = "cos"
```

- 备份失败告警通知，目前实现邮件告警

# 计划
- Web UI 操作页面
- 其他类型数据库数据备份，如 Mysql、Redis 等
- 探索非数据库类型的数据备份，如 Git 仓库项目代码、