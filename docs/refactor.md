# Go 重构进度

## 当前架构

项目采用轻量 DDD 分层：

- `cmd/big-market`：应用启动入口，完成配置加载、日志初始化、信号监听
- `internal/app`：应用装配层，统一完成依赖组装、HTTP 服务、定时任务、RabbitMQ 消费者生命周期管理
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

## 测试覆盖进度

- 策略规则模型解析：规则模型拆分、权重规则值解析、异常奖品 ID 过滤
- 策略 Redis 调度：普通概率表、权重概率表、未装配错误、库存扣减、售罄回写、脏数据解析错误
- 策略装配服务：策略库存缓存、默认概率表、权重概率表、活动映射、空奖品、缓存失败、规则缺失
- 策略库存同步：队列取值、空队列、落库更新、队列错误、仓储错误
- 策略查询服务：奖品列表锁规则、权重规则、空权重结果、仓储错误链路
- 责任链规则解析：黑名单规则、权重规则、非法规则格式
- 规则树引擎：库存接管、兜底奖品、节点跳转、缺失节点错误
- 活动参与下单：已有订单复用、活动状态、活动日期、总/月/日额度错误、已有月日账户复用
- 活动抽奖聚合：参与下单、策略抽奖、中奖记录保存、各步骤错误透传
- 活动账户与商品查询：参数校验、空账户兜底、日月账户兜底、仓储错误透传、商品查询错误
- 活动装配与库存：多 SKU 预热、仓储错误、缓存错误、空队列、库存售罄、库存归零消息
- 活动兑换：成功兑换、消息发送失败补偿、非法参数、活动状态、活动日期、库存不足
- 活动配送：参数裁剪、非法参数、仓储错误透传
- 返利处理：SKU 返利、积分返利、非法参数、未知返利类型、返利配置解析、活动状态错误
- 签到返利：签到订单、签到查询、空返利配置、仓储错误、保存错误、非法用户参数
- 积分账户：账户查询、空账户返回、非法用户参数、仓储错误
- 奖品服务：中奖记录保存、发奖消息、发布失败补偿、奖品分发、配置回查、配置错误、仓储错误
- RabbitMQ 消费者：发奖、返利、积分调账、活动库存归零的解析失败、业务错误、重复消费、订阅启动
- 消息任务补偿：待发送查询、发送成功、发送失败、状态更新失败、定时任务错误吞吐
- 基础设施保护：Redis 活动/策略 store 空连接、队列解析错误、RabbitMQ 发布/消费/关闭空连接
- 分库分表路由：Java mini-db-router 兼容、默认配置纠偏、空 key、单库默认路由
- MySQL 仓储查询与写入：活动基础信息、活动账户总/月/日、活动 SKU 与次数配置、活动参与订单查询、活动参与下单事务、活动参与新周期账户创建、活动兑换订单查询、活动兑换下单写入、活动配送边界与成功事务、活动配送账户创建、活动与策略库存写入、策略基础信息、策略奖品列表与明细、策略权重规则、策略锁规则、策略规则树装配、策略活动账户参与次数、奖品配置与待发奖任务、奖品发放事务、积分账户、积分支付完成与额度不足事务、积分返利入账事务、返利配置与签到订单、返利记录成对落库、积分兑换订单重复映射、任务状态更新 sqlmock 覆盖
- HTTP 接口：策略装配、策略奖品查询、策略权重查询、活动抽奖、活动账户、积分兑换、签到返利非法参数、health、ping、业务错误识别
- 通用类型与配置：响应包装、AppError、服务端口、MySQL、Redis、RabbitMQ、日志默认值

## 后续切片

- 继续补充活动仓储核心事务异常分支；优先覆盖 `SaveCreditPayOrder` 成功路径、月/日账户不存在时的参与下单路径、配送入账新建账户路径
- 拆分策略仓储中的规则树、规则权重、库存更新代码
- 补 HTTP 非法参数测试矩阵
- 补 RabbitMQ 消费者重试与失败日志细节测试
- 根据 Java 项目继续迁移缺失的运营、奖品、积分查询接口
