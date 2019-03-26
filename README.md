由于gin框架非常轻量，没有提供日常web开发常用的组件，所以本人又做了一层封装，方便基于gin的web程序开发。主要封装了mysql、redis、controller、conf、log、grpc、error、websocket等组件。
部分思路参考了之前公司老大封装的框架，特此声明。

### 运行例子
```
$ go get -u github.com/golang/protobuf/protoc-gen-go        #安装protoc-gen-go插件

#example提供http服务+grpc服务+websocket服务(使用ws-test-client.html测试)
$ cd example/proto
$ cp test.conf.default test.conf                            #按需修改配置
$ protoc --go_out=plugins=grpc:./ *.proto
$ go run main.go

#example2编写grpc客户端调用example提供的grpc服务
$ cd example2/proto
$ cp test.conf.default test.conf                            #按需修改配置
$ protoc --go_out=plugins=grpc:./ *.proto
$ go run main.go
```


### 配置说明
```
[system]
http_port=8080                      #http服务端口
rpc_port=8081                       #rpc服务端口
run_mode=debug                      #debug、test、release
worker_id=0                         #snowflake worker id，如，go项目：0-899，php项目：900-1023
param_secret=ceqcyxnprtj1t          #参数一致性秘钥

[log]
path=/Volumes/WorkHD/workspace/go/src/github.com/banbo/ys-gin/example/example.log
level=debug                         #debug、info、error

[db]
driver_name=mysql
host=127.0.0.1
port=3306
user=root
password=root
database=test
max_open=20
max_idle=10

[redis]
host=10.10.20.151
port=6379
password=123456
db=0

[rpc_client]
example_svr=localhost:8083          #rpc服务器地址
```