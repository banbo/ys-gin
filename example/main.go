package main

import (
	ysGin "github.com/banbo/ys-gin"

	"github.com/banbo/ys-gin/example/router"
	"github.com/banbo/ys-gin/example/rpc"
)

func main() {
	//启动服务
	app := ysGin.NewApp("/Volumes/WorkHD/workspace/go/src/github.com/banbo/ys-gin/example/test.conf")

	//设置路由
	router.Init(app.GinEngine)

	//注册rpc服务
	rpc.Register(app.RpcSvr)

	app.Run()
}
