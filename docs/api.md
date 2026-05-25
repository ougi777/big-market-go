# HTTP 接口说明

基础地址：

```text
http://localhost:8091
```

统一响应：

```json
{
  "code": "0000",
  "info": "success",
  "data": {}
}
```

## 健康检查

```http
GET /health
```

## 策略接口

装配抽奖策略：

```http
GET /api/v1/raffle/strategy/strategy_armory?strategyId=100001
```

随机抽奖：

```http
POST /api/v1/raffle/strategy/random_raffle
Content-Type: application/json

{"strategyId":100001}
```

查询活动奖品列表：

```http
POST /api/v1/raffle/strategy/query_raffle_award_list
Content-Type: application/json

{"userId":"xiaofuge","activityId":100301}
```

查询策略权重规则：

```http
POST /api/v1/raffle/strategy/query_raffle_strategy_rule_weight
Content-Type: application/json

{"userId":"xiaofuge","activityId":100301}
```

## 活动接口

装配活动库存和抽奖策略：

```http
GET /api/v1/raffle/activity/armory?activityId=100301
```

活动抽奖：

```http
POST /api/v1/raffle/activity/draw
Content-Type: application/json

{"userId":"xiaofuge","activityId":100301}
```

日历签到返利：

```http
POST /api/v1/raffle/activity/calendar_sign_rebate
Content-Type: application/x-www-form-urlencoded

userId=xiaofuge
```

查询当天是否已签到：

```http
POST /api/v1/raffle/activity/is_calendar_sign_rebate
Content-Type: application/x-www-form-urlencoded

userId=xiaofuge
```

查询用户活动账户：

```http
POST /api/v1/raffle/activity/query_user_activity_account
Content-Type: application/json

{"userId":"xiaofuge","activityId":100301}
```

查询活动 SKU：

```http
GET /api/v1/raffle/activity/query_sku_product_list_by_activity_id?activityId=100301
```

查询用户积分账户：

```http
GET /api/v1/raffle/activity/query_user_credit_account?userId=xiaofuge
```

积分兑换 SKU：

```http
POST /api/v1/raffle/activity/credit_pay_exchange_sku
Content-Type: application/json

{"userId":"xiaofuge","sku":9011}
```

## 常见错误码

```text
0002        参数错误
0003        唯一索引冲突
ERR_BIZ_001 策略 rule_weight 规则未配置
ERR_BIZ_002 抽奖策略未装配
ERR_BIZ_003 活动状态错误
ERR_BIZ_004 活动时间错误
ERR_BIZ_006 账户额度不足
ERR_BIZ_007 月账户额度不足
ERR_BIZ_008 日账户额度不足
ERR_BIZ_009 活动订单状态错误
```
