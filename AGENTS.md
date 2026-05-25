# bm-go 重构协作说明

## 项目定位

`bm-go` 是 `big-market-1-lee` 的 Go 版重构项目，目标是保留 Java 项目的业务模型和核心流程，用 Go 的工程习惯重建抽奖、活动、奖品、积分、返利、库存同步和消息补偿能力。

## 技术栈

- Golang
- Gin
- GORM
- MySQL
- Redis
- RabbitMQ
- robfig/cron
- Viper
- Zap

## 当前架构

```text
cmd/big-market              应用启动入口
configs                     本地配置
internal/config             配置加载
internal/trigger/http       HTTP 接口
internal/trigger/job        定时任务
internal/trigger/listener   RabbitMQ 消费者
internal/domain             领域模型、领域服务、仓储接口
internal/infrastructure     MySQL、Redis、RabbitMQ 实现
internal/types              通用响应、错误码、错误类型
```

## 已迁移能力

- 策略装配、Redis 概率表和抽奖调度。
- 抽奖前置责任链：黑名单、权重、默认规则。
- 抽奖后置规则树：解锁、库存、兜底奖品。
- 策略查询接口：奖品列表、权重规则。
- 活动装配、活动抽奖、活动账户查询。
- 活动 SKU 库存扣减、库存归零消费、库存同步任务。
- 奖品发奖消息、发奖消费、任务补偿。
- 日历签到返利、返利订单、返利消息消费。
- 用户积分账户查询、积分兑换 SKU、异步发货、消息补偿。
- 通用任务领域服务已从奖品服务拆出。
- 积分账户查询领域服务、积分交易仓储接口和积分 MySQL 仓储实现已从活动模块拆出。
- 用户相关分表表名路由，兼容 Java `mini-db-router` 表索引算法。
- MySQL 多数据源配置、运行时分库路由、消息任务多库补偿扫描。
- 策略错误码和 HTTP 业务错误透传。
- 发奖、兑换、返利和通用任务补偿的消息失败分支测试。

## 运行方式

```powershell
go test ./...
go build ./...
go vet ./...
go run ./cmd/big-market
```

健康检查：

```text
GET /health
```

## 重构原则

- 每次迁移一个清晰业务切片。
- 每个切片完成后运行 `go test ./...`、`go build ./...`、`go vet ./...`。
- 每个切片验证通过后单独提交，commit 使用中文。
- 优先沿用 Java 项目的业务语义和数据表结构。
- Go 代码保持轻量 DDD：领域包定义模型和接口，基础设施包实现 MySQL、Redis、RabbitMQ。
- 禁止批量删除文件或目录。

## 数据库说明

默认配置连接 `big_market` 单库，`sharding.db_count` 和 `sharding.table_count` 默认值为 `1`。访问 Java 分库分表数据时，可配置：

```yaml
mysql:
  shards:
    db01:
      dsn: root:123456@tcp(localhost:13308)/big_market_01?charset=utf8mb4&parseTime=True&loc=Local
    db02:
      dsn: root:123456@tcp(localhost:13308)/big_market_02?charset=utf8mb4&parseTime=True&loc=Local
sharding:
  db_count: 2
  table_count: 4
```

已支持表名分片：

```text
raffle_activity_order
user_raffle_order
user_award_record
user_behavior_rebate_order
user_credit_order
```

活动账户表、积分账户表和任务表按用户分库路由。

## 后续优先级

- 补齐抽奖、兑换、返利链路的端到端集成测试。
- 补齐 HTTP 接口文档和请求示例。
