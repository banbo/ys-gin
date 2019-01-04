ys-gin框架是基于gin框架开发。封装了mysql、redis、分布式id、controller、conf、log、grpc、error等组件。

example可以运行，封装了rpc服务端，example2里是rpc客户端。

### proto
```
$ cd example/proto
$ protoc --go_out=plugins=grpc:./ hello.proto  # *.proto
```
