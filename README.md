# gRPC Middleware

gRPC 拦截器主要分为两种：客户端拦截器（ClientInterceptor）服务端拦截器（ServerInterceptor），顾名思义，分别于请求的两端拦截 RPC 的执行。
一般用来实现熔断、限流、日志收集、open-tracing、异常捕获、数据统计、鉴权、数据注入等等多种功能。

### Usage

```shell
go get github.com/ZuoFuhong/grpc-middleware
```

**链路日志上报**

服务端添加拦截器，在服务被调用时，上报服务调用链路日志.

```go
import (
    "github.com/ZuoFuhong/grpc-middleware/tracing"
)

s := grpc.NewServer(grpc.UnaryInterceptor(tracing.UnaryServerInterceptor()))
```

客户端添加拦截器，在发起服务调用时，上报服务调用链路日志.

```go
import (
    "github.com/ZuoFuhong/grpc-middleware/tracing"
)

conn, err := grpc.DialContext(ctx, "127.0.0.1:1024", grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor()))
```

其中，全链路日志将通过 RPC 上报到 [grpc-datacollector](https://github.com/ZuoFuhong/grpc-datacollector) 集群中，详情参见架构文档。

**自定义编解码**

gRPC 提供了 codec 的插件注入能力，以实现自定义编解码.

服务调用时指定编码器，请求参数的 struct 对象可以序列化成 json.

```go
rpcRsp, err := stub.ImportWallet(context.Background(), &pb.ImportWalletReq{
    PrivateKey: "0x01c4bda0939df07a31e3738c6c1e1d5905c9f229e6ffa1922557308a62efb23f",
}, grpc.CallContentSubtype(codec.Name)) // 指定 JSON 编码
```

服务端需要隐式导包注册编码器，下游参数将能通过 json 反序列化成 struct 对象.

```go
import (
    _ "github.com/ZuoFuhong/grpc-middleware/encoding/json"
)
```

### License

This project is licensed under the [Apache 2.0 license](https://github.com/ZuoFuhong/grpc-middleware/blob/master/LICENSE).