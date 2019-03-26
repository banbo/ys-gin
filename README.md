- 由于gin框架非常轻量，没有提供日常web开发常用的组件，所以本人又做了一层封装，方便基于gin的web程序开发。主要封装了conf、log、orm、redis、grpc、error、controller、websocket等组件
- controller/controller主要封装了获取http参数方法、返回http数据(json)的方法
- controller/wsController主要封装了websocket返回数据方法
- errors目录里主要封装了自定义error接口，用于区分业务报错和系统报错
- middleware目录里主要封装了一些常用的中间件，如跨域、参数一致性验证等中间件
- test目录里主要封装了模拟浏览器请求，方便对http接口的单元测试，可参考example/controllers/test/testController_test.go
- conf基于beego/conf、log基于beego/log、orm基于xorm、redis基于go-redis等，这些都直接使用了非常成熟好用的组件，只是做了初始化封装，就没有重复造轮子了(貌似暂时也没这个实力^_^)
- 其余很多思路也都是参考各位大佬，不再一一列举，非常感谢大家的启发。

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