由于gin框架非常轻量，没有提供日常web开发常用的组件，所以本人又做了一层封装，方便基于gin的web程序开发。主要封装了conf、log、orm、redis、grpc、error、controller、websocket等组件。
conf、log、orm、redis这些都直接使用了非常成熟好用的组件，只是做了初始化封装，就没有重复造轮子了(貌似暂时也没这个实力^_^)。
其余很多思路也都是参考各位大佬，不再一一列举，非常感谢大家的启发。


### 目录结构
```
github.com
  - banbo
    - cache                         封装了缓存相关类，其中的redis基于go-redis
    - conf                          基于beego/conf
    - constant                      主要定义了一些公用的常量
    - controller
      - controller.go               封装了获取http参数的方法(不习惯gin的那套方法)、返回http数据(json)的方法
      - wsController.go             封装了返回websocket数据的方法
    - errors                        封装了自定义error接口，用于区分业务报错和系统报错
    - id                            封装了分布式的snowflak id
    - log                           基于beego/log
    - middleware                    封装了一些常用的中间件，如跨域、参数一致性验证等中间件
    - model                         封装了orm和分页方法，orm是基于xorm
    - rpc                           封装了响应错误数据的方法
    - test                          封装了模拟浏览器请求，方便对http接口的单元测试，可参考example/controllers/test/testController_test.go
    - utils                         封装了常用的工具类、工具函数
```


### 配置说明
```
[system]
http_port=8080                      #http服务端口
rpc_port=8081                       #rpc服务端口
run_mode=debug                      #debug、test、release
worker_id=0                         #机器id，用于生成SnowflakeID，go项目：0-899，php项目：900-1023
param_secret=ceqcyxnprtj1t          #参数一致性秘钥
dbs=mysql,sqlite3                   #多数据库实例，具体配置见下db-mysql、db-sqlite3

[log]
path=/Volumes/WorkHD/workspace/go/src/github.com/banbo/ys-gin/example/example.log
level=debug                         #debug、info、error

[db-mysql]
driver_name=mysql
host=127.0.0.1
port=3306
user=root
password=root
database=test
max_open=20                         #最大连接数
max_idle=10                         #最大空闲连接数

[db-sqlite3]
driver_name=sqlite3
database=./db/main

[redis]
host=10.10.20.151
port=6379
password=123456
db=0

[rpc_client]
example_svr=localhost:8083          #rpc服务器地址
```


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
