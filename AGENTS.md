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

- 策略装配、Redis 概率表和抽奖调度
- 抽奖前置责任链：黑名单、权重、默认规则
- 抽奖后置规则树：解锁、库存、兜底奖品
- 策略查询接口：奖品列表、权重规则
- 活动装配、活动抽奖、活动账户查询
- 活动 SKU 库存扣减、库存归零消费、库存同步任务
- 奖品发奖消息、发奖消费、任务补偿
- 日历签到返利、返利消息消费
- 用户积分账户查询、积分兑换 SKU、异步发货、消息补偿
- 用户相关分表表名路由，默认 `table_count: 1`

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

- 每次只迁移一个清晰业务切片。
- 每个切片完成后运行 `go test ./...`、`go build ./...`、`go vet ./...`。
- 每个切片验证通过后单独提交，commit 使用中文。
- 优先沿用 Java 项目的业务语义和数据表结构。
- Go 代码保持轻量 DDD：领域包定义模型和接口，基础设施包实现 MySQL、Redis、RabbitMQ。
- 禁止批量删除文件或目录。

## 数据库说明

默认配置连接 `big_market` 单库，`sharding.table_count` 默认为 `1`。需要访问 Java 分表数据时，可设置：

```yaml
sharding:
  table_count: 4
```

当前已支持用户流水类表的表名分片：`raffle_activity_order`、`user_raffle_order`、`user_award_record`、`user_behavior_rebate_order`、`user_credit_order`。多数据源库路由仍需后续切片实现。

## 后续优先级

- 多数据源库路由，兼容 `big_market_01`、`big_market_02`
- 分表路由在活动账户日/月表、任务表上的完整对齐
- 积分领域从活动仓储中拆出独立仓储
- 抽奖、兑换、返利链路的集成测试
- README 乱码修复和接口文档补全
