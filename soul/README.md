## gRPC + REST Gateway Play

Blog post: https://coreos.com/blog/gRPC-protobufs-swagger.html

To try it all out do this:

```
$ go get -u github.com/philips/grpc-gateway-example
$ grpc-gateway-example serve
$ grpc-gateway-example echo "my first rpc echo"
$ curl -X POST -k https://localhost:10000/v1/echo -H "Content-Type: text/plain" -d '{"value": "foo"}'
{"value":"my REST echo"}
```


Huge thanks to the hard work people have put into the [Go gRPC bindings][gogrpc] and [gRPC to JSON Gateway][grpcgateway]

[gogrpc]: https://github.com/grpc/grpc-go
[grpcgateway]: https://github.com/grpc-ecosystem/grpc-gateway
# soul
一、准备工作
1、安装protobuf
https://blog.csdn.net/wwwyuanliang10000/article/details/78923137
2. 安装gRPC-go
安装golang protobuf直接使用golang的get即可
go get -u github.com/golang/protobuf/proto // golang protobuf 库
go get -u github.com/golang/protobuf/protoc-gen-go //protoc –go_out 工具
3、获取grpc
go get google.golang.org/grpc