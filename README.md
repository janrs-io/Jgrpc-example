# Jgrpc-example

基于https://github.com/janrs-io/Jgrpc 微服务最佳实践写的几个案例。

**仅提供一个 `Demo` 案例代码**

*本项目需要使用 `k8s` 集群才能跑通。如果是本地启动服务的话，需要把每个服务的客户端以及服务端的端口都改成不一样的。*

# 已开发功能

## `auth` 服务

- 注册授权功能
- 删除授权功能
- 刷新授权功能
- 获取授权功能
- 接口白名单功能

## `user` 服务

- 用户注册功能
- 用户登录功能
- 用户退出功能
- 获取用户详情功能
- 获取用户订单详情[该接口会调用user-grpc/order-grpc/product-grpc接口，并且使用dtm事务以及绕过auth白名单]

## `product` 产品服务

- 添加产品
- 删除产品
- 获取产品详情
- 获取产品列表

## `order` 服务

- 添加订单功能 [增加了 DTM 事务支持。创建订单会使用 DTM 的 SAGA 事务扣除产品库存]
- 获取订单详情功能