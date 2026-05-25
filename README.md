# bm-go

`bm-go` 是 `big-market-1-lee` 的 Go 版重构项目。当前工程按轻量 DDD 分层组织，逐步迁移 Java 项目的策略抽奖、活动账户、奖品发放、积分交易、签到返利、库存同步和消息补偿能力。

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

## 目录结构

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

## 当前进度

- 策略模块：策略装配、责任链、规则树、Redis 概率表、库存同步。
- 活动模块：活动装配、活动抽奖、活动账户查询、SKU 查询、积分兑换。
- 奖品模块：中奖记录、发奖消息、发奖消费、任务补偿。
- 返利模块：日历签到、返利订单、返利消息消费。
- 积分模块：积分账户查询、积分扣减、异步发货消息。
- 基础设施：MySQL、Redis、RabbitMQ、Cron、用户分表表名路由。

## 本地运行

```powershell
go mod tidy
go test ./...
go build ./...
go vet ./...
go run ./cmd/big-market
```

默认配置文件：

```text
configs/config.yaml
```

健康检查：

```text
GET /health
```

## 策略接口

```text
GET  /api/v1/raffle/strategy/strategy_armory?strategyId=100001
POST /api/v1/raffle/strategy/random_raffle
POST /api/v1/raffle/strategy/query_raffle_award_list
POST /api/v1/raffle/strategy/query_raffle_strategy_rule_weight
```

## 活动接口

```text
GET  /api/v1/raffle/activity/armory?activityId=100301
POST /api/v1/raffle/activity/draw
POST /api/v1/raffle/activity/calendar_sign_rebate
POST /api/v1/raffle/activity/is_calendar_sign_rebate
POST /api/v1/raffle/activity/query_user_activity_account
GET  /api/v1/raffle/activity/query_sku_product_list_by_activity_id?activityId=100301
GET  /api/v1/raffle/activity/query_user_credit_account?userId=xiaofuge
POST /api/v1/raffle/activity/credit_pay_exchange_sku
```

## 分表配置

默认使用单表：

```yaml
sharding:
  table_count: 1
```

需要访问 Java 项目的分表数据时，可调整为：

```yaml
sharding:
  table_count: 4
```

当前支持分表表名路由的用户流水表：

```text
raffle_activity_order
user_raffle_order
user_award_record
user_behavior_rebate_order
user_credit_order
```

多数据源库路由仍在后续重构范围内。
