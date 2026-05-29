# 项目改进总览

## 文档目标

这份文档记录 `bm-go` 后续改进方向，帮助理解每个改进背后的业务功能、技术原因和执行顺序。当前讨论的 RabbitMQ 生产者 confirm 属于“消息可靠性”改进的一部分。

## 项目业务主线

`bm-go` 的核心业务是围绕活动抽奖和积分权益完成一组最终一致流程：

1. 用户参加活动，系统扣减活动账户次数。
2. 抽奖策略给出中奖结果，系统保存中奖记录。
3. 发奖消息进入 RabbitMQ，消费者完成奖品发放。
4. 用户可通过积分兑换活动 SKU，系统扣减积分并发放活动额度。
5. 签到返利、积分返利、库存归零等动作通过消息和定时任务补齐异步流程。

这类业务的关键点是“先落业务事实，再异步推进后续动作”。数据库记录表达业务事实，RabbitMQ 推动异步动作，任务表负责消息补偿。

## 改进总览

| 优先级 | 改进方向 | 业务价值 | 涉及模块 |
| --- | --- | --- | --- |
| P0 | RabbitMQ 生产者可靠投递 | 提升发奖、返利、积分调账消息到达 Broker 的确定性 | `internal/infrastructure/rabbitmq`、`internal/domain/task` |
| P0 | 消息任务补偿统一 | 让所有关键异步消息拥有一致的失败恢复路径 | `internal/domain/task`、活动库存、发奖、返利、兑换 |
| P1 | 消费端重试与失败记录 | 降低业务处理失败后的排查成本，增强消费恢复能力 | `internal/trigger/listener` |
| P1 | 端到端集成测试 | 验证抽奖、兑换、返利链路在真实依赖下的完整行为 | `tests`、`internal/app` |
| P1 | HTTP 接口文档和示例 | 支撑联调、压测、面试讲解和项目展示 | `docs/api.md` |
| P2 | 可观测性增强 | 快速定位消息积压、任务失败、库存同步异常 | 日志、指标、追踪 |
| P2 | 配置治理 | 提高本地、测试、生产环境切换效率 | `configs`、`internal/config` |

## P0：RabbitMQ 生产者可靠投递

### 业务功能

RabbitMQ 负责承接这些关键异步动作：

- `send_award`：中奖记录保存后，异步发放奖品。
- `send_rebate`：签到返利订单保存后，异步处理 SKU 或积分返利。
- `credit_adjust_success`：积分扣减成功后，异步发放活动额度。
- `activity_sku_stock_zero`：活动 SKU 库存归零后，异步清理库存队列。

这些消息连接了“业务记录已落库”和“后续动作已执行”之间的空档。

### 当前表现

[client.go](D:/StudyCode/aiCode/bm-go/internal/infrastructure/rabbitmq/client.go:36) 当前发布流程是：

```text
创建 channel
声明 durable queue
PublishWithContext 发送消息
返回 Publish 调用结果
```

这个流程能发现连接错误、channel 错误和基础发布错误。Broker 接收后的确认结果需要 producer confirm 补齐。

### 改进原因

业务上，发奖、返利、积分调账都属于关键消息。生产者 confirm 可以让发布方知道 Broker 对消息的处理结果：

- `ack`：Broker 已确认接收。
- `nack`：Broker 明确拒绝或处理失败。
- `timeout`：发布方在指定时间内没有收到确认。

confirm 配合本地任务表，可以形成更完整的可靠投递链路：

```text
业务事务写入 task(create)
发布 MQ 并等待 confirm
confirm ack -> task=completed
confirm nack/timeout/error -> 有限重试
重试耗尽 -> task=fail
定时补偿扫描 create/fail -> 重新发布
```

### 执行计划

1. 在 `Publish` 中开启 confirm。
2. 使用 `NotifyPublish` 接收 `ack/nack`。
3. 发布消息时设置 `DeliveryMode: amqp.Persistent`。
4. 收到 `nack`、confirm 超时、连接短暂异常时执行有限次即时重试。
5. 重试耗尽后返回错误，由领域服务把 task 标记为 `fail`。
6. 保持 `Publish(ctx, topic, message)` 方法签名稳定，业务层继续使用当前接口。

建议默认参数：

```text
重试次数：3
退避间隔：50ms、100ms、200ms
confirm 等待：跟随 ctx deadline
重试粒度：每次重新创建 channel
```

## P0：消息任务补偿统一

### 业务功能

任务表是本项目异步消息的本地消息表。业务事务保存核心数据时，同时保存待发送消息。消息发布成功后任务变为 `completed`，发布失败后任务变为 `fail`，定时任务继续扫描并补偿。

### 当前表现

这些链路已经接入任务补偿：

- 发奖消息：中奖记录和发奖 task 同步保存。
- 返利消息：返利订单和返利 task 同步保存。
- 积分兑换成功消息：积分流水、兑换订单和成功通知 task 同步保存。

活动库存归零消息当前由库存扣减逻辑直接发布。它适合纳入任务表，让库存归零清理也拥有统一补偿路径。

### 改进原因

统一补偿模型可以让所有关键消息遵循相同规则：

```text
先写业务事实
再发消息
发送失败进入 fail
定时任务持续补偿
消费者幂等处理重复消息
```

这条规则便于开发、排查和面试表达。

### 执行计划

1. 梳理所有 `Publish` 调用点。
2. 标记关键业务消息和普通通知消息。
3. 将关键业务消息全部绑定 task 表。
4. 让补偿任务扫描 `create/fail` 状态。
5. 用 `message_id`、业务单号、唯一索引保证消费幂等。

## P1：消费端重试与失败记录

### 业务功能

消费者负责把消息转换成业务动作。例如发奖消费者把 `send_award` 转成奖品发放，返利消费者把 `send_rebate` 转成积分入账或活动额度发放。

### 改进原因

当前消费者处理失败会 `Nack(false, true)`，消息重新入队。后续可以加入失败次数、错误日志和死信队列，形成更清晰的故障处理路径。

### 执行计划

1. 给消费者增加统一错误分类。
2. 对可恢复错误执行有限重试。
3. 对解析错误、长期业务错误进入失败记录或死信队列。
4. 日志记录 `topic`、`message_id`、`user_id`、业务单号、错误原因。

## P1：端到端集成测试

### 业务功能

端到端测试验证从 HTTP 请求到数据库、Redis、RabbitMQ、消费者的完整流程。

### 改进原因

当前单元测试覆盖了大量领域逻辑。端到端测试可以证明核心业务链路在真实依赖组合下运行稳定。

### 执行计划

1. 抽奖链路：活动装配、参与活动、策略抽奖、中奖记录、发奖消息。
2. 积分兑换链路：创建兑换订单、扣减积分、发送成功消息、发放活动额度。
3. 签到返利链路：签到、返利订单、返利消息、积分或额度入账。
4. 消息补偿链路：模拟发布失败，验证 task 从 `fail` 到 `completed`。

## P1：HTTP 接口文档和请求示例

### 业务功能

HTTP 接口是项目外部入口。清晰的接口文档可以支撑自测、联调和项目讲解。

### 执行计划

1. 为每个接口补齐请求参数、响应示例和业务错误码。
2. 增加典型业务流程调用顺序。
3. 增加 Postman 或 curl 示例。
4. 标注接口背后的领域服务和数据变化。

## P2：可观测性增强

### 业务功能

可观测性帮助定位消息积压、补偿失败、库存同步延迟、接口错误。

### 执行计划

1. 关键日志增加 `message_id`、`topic`、`user_id`、`order_id`。
2. 定时任务输出扫描数量、成功数量、失败数量。
3. RabbitMQ 发布和消费记录耗时。
4. 后续接入指标面板：消息发布成功率、补偿成功率、消费者失败率。

## P2：配置治理

### 业务功能

配置负责区分本地开发、测试环境和生产环境。RabbitMQ、MySQL、Redis、任务周期都依赖配置。

### 执行计划

1. 拆分 `config.local.yaml`、`config.test.yaml`、`config.prod.yaml` 示例。
2. RabbitMQ 增加 publish retry、confirm timeout、consumer prefetch 等配置。
3. 定时任务周期按场景配置。
4. 文档说明环境变量覆盖规则。

## 执行记录

### 2026-05-29：RabbitMQ 生产者可靠投递

已完成第一阶段改进：

1. `Publish` 开启 producer confirm。
2. 发布前声明 durable queue。
3. 发布消息设置 `DeliveryMode: amqp.Persistent`。
4. 发布后等待 `ack/nack`。
5. 收到 `nack`、confirm 超时、连接或 channel 发布错误时执行有限重试。
6. 重试耗尽后返回错误，业务层沿用现有 task 标记 `fail` 和定时补偿。
7. 方法签名保持 `Publish(ctx, topic, message)`，领域层无需调整。

已补充测试：

- confirm ack 发布成功。
- confirm nack 后重试成功。
- confirm nack 重试耗尽返回错误。
- confirm 等待超时返回错误。
- 空连接发布、消费、关闭保护保持通过。

已验证：

```powershell
go test ./...
go build ./...
go vet ./...
```

## 当前推荐执行顺序

1. 已完成 RabbitMQ confirm、消息持久化、有限重试。
2. 继续验证现有 task 补偿逻辑，重点关注发奖、返利、兑换链路。
3. 将活动库存归零消息纳入 task 补偿。
4. 增强消费者失败日志和重试策略。
5. 补齐端到端集成测试。
6. 完善 HTTP 接口文档和请求示例。

最终目标是形成一条清晰可靠的异步业务链路：数据库保存业务事实，RabbitMQ 推动异步动作，confirm 提供生产者确认，task 表提供最终补偿，消费者幂等保证重复消息安全。
