# Go 重构进度

## 当前架构

项目采用轻量 DDD 分层：

- `cmd/big-market`：应用启动入口，完成配置加载、依赖装配、HTTP 与定时任务启动
- `internal/domain`：领域模型、仓储接口、领域服务、规则引擎
- `internal/infrastructure`：MySQL、Redis、RabbitMQ 等基础设施实现
- `internal/trigger/http`：Gin HTTP 接口适配
- `internal/trigger/job`：库存同步、消息补偿等定时任务
- `internal/trigger/listener`：RabbitMQ 消费者

## 已迁移能力

- 策略装配、抽奖责任链、规则树过滤、Redis 概率表调度
- 活动装配、活动抽奖、活动账户查询、SKU 商品查询
- 积分兑换 SKU、活动订单配送、账户额度增加与扣减
- 奖品发放、中奖记录、消息任务补偿
- 日历签到返利、返利订单、积分账户发放
- 分库分表路由、任务表多库扫描

## 活动仓储拆分

活动仓储按能力拆分为多个文件，仍由 `ActivityRepository` 统一实现领域接口：

- `activity_repository.go`：结构体、构造函数、基础活动查询
- `activity_account_repository.go`：活动账户总/月/日查询
- `activity_account_quota_repository.go`：账户额度扣减、额度补充、月日镜像更新
- `activity_partake_repository.go`：参与抽奖订单查询与保存
- `activity_sku_repository.go`：SKU 商品与活动次数配置查询
- `activity_exchange_repository.go`：积分兑换订单查询与保存
- `activity_rebate_repository.go`：返利 SKU 订单保存
- `activity_delivery_repository.go`：活动订单配送与额度入账
- `activity_stock_repository.go`：活动 SKU 库存落库更新
- `activity_task_repository.go`：活动消息任务状态更新

## 策略仓储拆分

策略仓储按抽奖策略核心能力拆分为多个文件，仍由 `StrategyRepository` 统一实现领域接口：

- `strategy_repository.go`：结构体、构造函数、DB helper、接口断言
- `strategy_base_repository.go`：策略基础查询、活动 ID 到策略 ID 映射
- `strategy_award_repository.go`：策略奖品列表与单个奖品查询
- `strategy_rule_repository.go`：策略规则、奖品规则模型、规则值、锁规则统计
- `strategy_rule_tree_repository.go`：规则树、规则节点、规则连线装配查询
- `strategy_rule_weight_repository.go`：权重规则解析与权重奖品查询
- `strategy_activity_account_repository.go`：活动账户参与次数、日参与次数、今日抽奖次数查询
- `strategy_stock_repository.go`：策略奖品库存扣减队列与库存落库更新
- `strategy_dispatch.go`：Redis 概率表随机调度

## 验证命令

每个重构切片完成后执行：

```powershell
go test ./...
go build ./...
go vet ./...
```

## 后续切片

- 给活动仓储核心事务增加 sqlmock 或集成测试
- 拆分策略仓储中的规则树、规则权重、库存更新代码
- 补 HTTP 非法参数测试矩阵
- 补 RabbitMQ 消费者重试与失败日志测试
- 根据 Java 项目继续迁移缺失的运营、奖品、积分查询接口
