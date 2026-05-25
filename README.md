# bm-go

`bm-go` 是 `big-market` Java 项目的 Go 重构骨架。当前目录按照原项目的 DDD 分层做 Go 化拆分，后续逐步迁移抽奖策略、活动、奖品、积分、返利和任务补偿逻辑。

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
configs                     本地配置模板
internal/app                运行期装配根
internal/config             配置加载
internal/trigger/http       HTTP 入口
internal/trigger/job        定时任务入口
internal/trigger/listener   MQ 消费入口
internal/domain             领域模型和领域服务
internal/infrastructure     MySQL、Redis、RabbitMQ 实现
internal/types              通用响应、错误码和错误类型
```

## 本地启动

```powershell
go mod tidy
go run ./cmd/big-market
```

健康检查：

```text
GET /health
```

策略接口：

```text
GET  /api/v1/raffle/strategy/strategy_armory?strategyId=100001
POST /api/v1/raffle/strategy/random_raffle
POST /api/v1/raffle/strategy/query_raffle_award_list
POST /api/v1/raffle/strategy/query_raffle_strategy_rule_weight
```
